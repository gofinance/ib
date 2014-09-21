// Package trade offers a pure Go abstraction over Interactive Brokers IB API.
//
// Engine is the main type. It provides a mechanism to connect to either IB
// Gateway or TWS, send Request values and receive Reply values. The Engine
// provides an observer pattern both for receiving Reply values as well as Engine
// termination notification. Any network level errors will terminate the Engine.
//
// A high-level Manager interface is also provided. This provides a way to
// easily use IB API without needing to deal directly with Engine and the
// associated Request, Reply, message ID and Reply ordering considerations.
//
// All types are thread-safe and can be used from multiple goroutines at once.
// Blocking methods are identified in the documentation.
package ib

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"time"
)

const (
	gatewayDefault   = "127.0.0.1:4001"
	UnmatchedReplyId = int64(-9223372036854775808)
)

type NewEngineOptions struct {
	Gateway          string
	Client           int64
	DumpConversation bool
}

// Engine is the entry point to the IB IB API
type Engine struct {
	id               chan int64
	exit             chan bool
	terminated       chan struct{}
	ch               chan command
	gateway          string
	client           int64
	con              net.Conn
	reader           *bufio.Reader
	input            *bytes.Buffer
	output           *bytes.Buffer
	rxReply          chan Reply
	rxErr            chan error
	txRequest        chan txrequest
	txErr            chan error
	observers        map[int64]chan<- Reply
	unObservers      []chan<- Reply
	allObservers     []chan<- Reply
	stObservers      []chan<- EngineState
	state            EngineState
	serverTime       time.Time
	clientVersion    int64
	serverVersion    int64
	dumpConversation bool
	lastDumpRead     int64
	lastDumpId       int64
	fatalError       error
}

type command struct {
	fun func()
	ack chan struct{}
}

type txrequest struct {
	req Request
	ack chan struct{}
}

func uniqueId(start int64) chan int64 {
	ch := make(chan int64)
	id := start
	go func() {
		for {
			if id == UnmatchedReplyId {
				id += 1
			}
			ch <- id
			id += 1
		}
	}()
	return ch
}

// Next client id. Package scope to ensure new engines have unique client IDs.
var clientSeq = uniqueId(1)

// NewEngine takes a client id and returns a new connection
// to IB Gateway or IB Trader Workstation.
func NewEngine(opt NewEngineOptions) (*Engine, error) {
	gateway := opt.Gateway
	if gateway == "" {
		gateway = gatewayDefault
	}
	con, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}

	client := opt.Client
	if client == 0 {
		client = <-clientSeq
	}

	e := Engine{
		id:               uniqueId(100),
		exit:             make(chan bool),
		terminated:       make(chan struct{}),
		ch:               make(chan command),
		gateway:          gateway,
		client:           client,
		con:              con,
		reader:           bufio.NewReader(con),
		input:            bytes.NewBuffer(make([]byte, 0, 4096)),
		output:           bytes.NewBuffer(make([]byte, 0, 4096)),
		rxReply:          make(chan Reply),
		rxErr:            make(chan error),
		txRequest:        make(chan txrequest),
		txErr:            make(chan error),
		observers:        make(map[int64]chan<- Reply),
		state:            EngineReady,
		dumpConversation: opt.DumpConversation,
	}

	err = e.handshake()
	if err != nil {
		return nil, err
	}

	// start worker goroutines (these exit on request or error)
	go e.startReceiver()
	go e.startTransmitter()
	go e.startMainLoop()

	return &e, nil
}

func (e *Engine) handshake() error {
	// write client version and id
	clientShake := &clientHandshake{clientVersion, e.client}
	e.output.Reset()
	if err := clientShake.write(e.output); err != nil {
		return err
	}

	if _, err := e.con.Write(e.output.Bytes()); err != nil {
		return err
	}

	// read server version and time
	serverShake := &serverHandshake{}
	e.input.Reset()
	if err := serverShake.read(e.reader); err != nil {
		return err
	}

	if serverShake.version < minServerVersion {
		return fmt.Errorf("%s must be at least version %d (reported %d)",
			e.ConnectionInfo(), minServerVersion, serverShake.version)
	}

	e.serverVersion = serverShake.version
	e.serverTime = serverShake.time
	return nil
}

func (e *Engine) startReceiver() {
	defer func() {
		close(e.rxReply)
		close(e.rxErr)
	}()
	for {
		r, err := e.receive()
		if err != nil {
			select {
			case <-e.terminated:
				return
			case e.rxErr <- err:
				return
			}
		}

		select {
		case <-e.terminated:
			return
		case e.rxReply <- r:
		}
	}
}

func (e *Engine) startTransmitter() {
	defer func() {
		// Don't close txRequest, as we are not the sender
		close(e.txErr)
	}()
	for {
		select {
		case <-e.terminated:
			return
		case t := <-e.txRequest:
			err := e.transmit(t.req)
			if err != nil {
				select {
				case <-e.terminated:
					return
				case e.txErr <- err:
					return
				}
			}
			close(t.ack)
		}
	}
}

func (e *Engine) startMainLoop() {
	defer func() {
		// Signal terminating for benefit of other goroutines
		close(e.terminated)

		// Safe to kill the connection, as we're advising other goroutines we're quitting
		e.con.Close()

		// Wait for other goroutines to indicate they've finished
		<-e.txErr
		<-e.rxErr

	outer:
		for _, ob := range e.stObservers {
			for {
				select {
				case ob <- e.state:
					continue outer
				case <-time.After(time.Duration(5) * time.Second):
					log.Printf("Waited 5 seconds for state channel %v\n", ob)
				}
			}
		}
	}()
	for {
		select {
		case <-e.exit:
			e.state = EngineExitNormal
			return
		case err := <-e.rxErr:
			log.Printf("%s engine: RX error %s", e.ConnectionInfo(), err)
			e.fatalError = err
			e.state = EngineExitError
			return
		case err := <-e.txErr:
			log.Printf("%s engine: TX error %s", e.ConnectionInfo(), err)
			e.fatalError = err
			e.state = EngineExitError
			return
		case cmd := <-e.ch:
			cmd.fun()
			close(cmd.ack)
		case r := <-e.rxReply:
			e.deliverToObservers(r)
		}
	}
}

func (e *Engine) deliverToObservers(r Reply) {
	if r.code() == mErrorMessage {
		var done []chan<- Reply
		for _, o := range e.observers {
			for _, prevDone := range done {
				if o == prevDone {
					continue
				}
			}
			done = append(done, o)
			e.deliverToObserver(o, r)
		}
		for _, o := range e.unObservers {
			e.deliverToObserver(o, r)
		}
		for _, o := range e.allObservers {
			e.deliverToObserver(o, r)
		}
		return
	}
	if mr, ok := r.(MatchedReply); ok {
		if o, ok := e.observers[mr.Id()]; ok {
			e.deliverToObserver(o, r)
		}

		for _, o := range e.allObservers {
			e.deliverToObserver(o, r)
		}
		return
	}
	// must be a non-error, unmatched reply
	for _, o := range e.unObservers {
		e.deliverToObserver(o, r)
	}
	for _, o := range e.allObservers {
		e.deliverToObserver(o, r)
	}
}

func (e *Engine) deliverToObserver(c chan<- Reply, r Reply) {
	for {
		select {
		case c <- r:
			return
		case <-time.After(time.Duration(5) * time.Second):
			log.Printf("Waited 5 seconds for reply channel %v\n", c)
		}
	}
}

func (e *Engine) transmit(r Request) (err error) {
	e.output.Reset()

	// encode message type and client version
	hdr := &header{
		code:    int64(r.code()),
		version: r.version(),
	}

	if err = hdr.write(e.output); err != nil {
		return
	}

	// encode the message itself
	if err = r.write(e.output); err != nil {
		return
	}

	if e.dumpConversation {
		b := e.output
		s := strings.Replace(b.String(), "\000", "-", -1)
		fmt.Printf("%d> '%s'\n", e.client, s)
	}

	_, err = e.con.Write(e.output.Bytes())
	return
}

// NextRequestId returns a unique request id (which is never UnmatchedReplyId).
func (e *Engine) NextRequestId() int64 {
	return <-e.id
}

func (e *Engine) ClientId() int64 {
	return e.client
}

// ConnectionInfo returns the gateway address and client ID of this connection.
func (e *Engine) ConnectionInfo() string {
	return fmt.Sprintf("%s/%d", e.gateway, e.client)
}

// sendCommand delivers the func to the engine, blocking the calling goroutine
// until the command is acknowledged as completed or the engine exits.
func (e *Engine) sendCommand(c func()) {
	cmd := command{c, make(chan struct{})}

	// send cmd
	select {
	case <-e.terminated:
		return
	case e.ch <- cmd:
	}

	// await ack (also handle termination, although it shouldn't happen
	// given the cmd was delivered so we beat any exit/error situations)
	select {
	case <-e.terminated:
		log.Println("Engine unexpectedly terminated after command sent")
		return
	case <-cmd.ack:
		return
	}
}

// Subscribe will notify subscribers of future events with given id.
// Many request types implement MatchedRequest and therefore provide a SetId().
// To receive the corresponding MatchedReply events, firstly subscribe with the
// same id as will be assigned with SetId(). Any incoming events that do not
// implement MatchedReply will be delivered to those observers subscribed to
// the UnmatchedReplyId constant. Note that the engine will raise an error if
// an attempt is made to send a MatchedRequest with UnmatchedReplyId as its id,
// given the high unlikelihood of that id being required in normal situations
// and that NextRequestId() guarantees to never return UnmatchedReplyId.
// Each ErrorMessage event is delivered once only to each known observer.
// The engine never closes the channel (allowing reuse across IDs and engines).
// This call will block until the subscriber is registered or engine terminates.
func (e *Engine) Subscribe(o chan<- Reply, id int64) {
	e.sendCommand(func() {
		if id != UnmatchedReplyId {
			e.observers[id] = o
			return
		}

		e.unObservers = append(e.unObservers, o)
	})
}

func (e *Engine) SubscribeAll(o chan<- Reply) {
	e.sendCommand(func() {
		e.allObservers = append(e.allObservers, o)
	})
}

// Unsubscribe blocks until the observer is removed. It also maintains a
// goroutine to sink the channel until the unsubscribe is finalised, which
// frees the caller from maintaining a separate goroutine.
func (e *Engine) Unsubscribe(o chan Reply, id int64) {
	terminate := make(chan struct{})
	go func() {
		for {
			select {
			case <-o:
			case <-terminate:
				return
			}
		}
	}()
	e.sendCommand(func() {
		if id != UnmatchedReplyId {
			delete(e.observers, id)
			return
		}

		newUnObs := make([]chan<- Reply, 0)
		for _, existing := range e.unObservers {
			if existing != o {
				newUnObs = append(newUnObs, o)
			}
		}
		e.unObservers = newUnObs
	})
	close(terminate)
}

func (e *Engine) UnsubscribeAll(o chan Reply) {
	terminate := make(chan struct{})
	go func() {
		for {
			select {
			case <-o:
			case <-terminate:
				return
			}
		}
	}()
	e.sendCommand(func() {
		newUnObs := make([]chan<- Reply, 0)
		for _, existing := range e.allObservers {
			if existing != o {
				newUnObs = append(newUnObs, o)
			}
		}
		e.allObservers = newUnObs
	})
	close(terminate)
}

// SubscribeState will register an engine state subscriber that is notified when
// the engine exits for any reason. The engine will close the channel after use.
// This call will block until the subscriber is registered or engine terminates.
func (e *Engine) SubscribeState(o chan<- EngineState) {
	if o == nil {
		return
	}
	e.sendCommand(func() { e.stObservers = append(e.stObservers, o) })
}

// UnsubscribeState blocks until the observer is removed. It also maintains a
// goroutine to sink the channel until the unsubscribe is finalised, which
// frees the caller from maintaining a separate goroutine.
func (e *Engine) UnsubscribeState(o chan EngineState) {
	terminate := make(chan struct{})
	go func() {
		for {
			select {
			case <-o:
			case <-terminate:
				return
			}
		}
	}()
	e.sendCommand(func() {
		var r []chan<- EngineState
		for _, exist := range e.stObservers {
			if exist != o {
				r = append(r, exist)
			}
		}
		e.stObservers = r
	})
	close(terminate)
}

// FatalError returns the error which caused termination (or nil if no error).
func (e *Engine) FatalError() error {
	return e.fatalError
}

// State returns the engine's current state.
func (e *Engine) State() EngineState {
	return e.state
}

// Stop blocks until the engine is fully stopped. It can be safely called on an
// already-stopped or stopping engine.
func (e *Engine) Stop() {
	select {
	case <-e.terminated:
		return
	case e.exit <- true:
	}

	<-e.terminated
}

type header struct {
	code    int64
	version int64
}

func (v *header) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.code); err != nil {
		return
	}
	return writeInt(b, v.version)
}

func (v *header) read(b *bufio.Reader) (err error) {
	if v.code, err = readInt(b); err != nil {
		return
	}
	v.version, err = readInt(b)
	return
}

// Send a message to the engine, blocking until sent or the engine exits.
// This method will return an error if the UnmatchedReplyId is used or the
// engine exits. A nil error indicates successful transmission. Any transmission
// failure (eg connectivity loss) will cause the engine to exit with an error.
func (e *Engine) Send(r Request) (err error) {
	if mr, ok := r.(MatchedRequest); ok {
		if mr.Id() == UnmatchedReplyId {
			return fmt.Errorf("%d is a reserved ID (try using NextRequestId)", UnmatchedReplyId)
		}
	}
	t := txrequest{r, make(chan struct{})}

	// send tx request
	select {
	case <-e.terminated:
		err = e.FatalError()
		if err == nil {
			err = fmt.Errorf("Engine has already exited normally")
		}
		return err
	case e.txRequest <- t:
	}

	// await ack or error
	select {
	case <-e.terminated:
		err = e.FatalError()
		if err == nil {
			err = fmt.Errorf("Engine has already exited normally")
		}
		return err
	case <-t.ack:
		return nil
	}
}

type packetError struct {
	value interface{}
	kind  reflect.Type
}

func (e *packetError) Error() string {
	return fmt.Sprintf("don't understand packet '%v' of type '%v'",
		e.value, e.kind)
}

func failPacket(v interface{}) error {
	return &packetError{
		value: v,
		kind:  reflect.ValueOf(v).Type(),
	}
}

func (e *Engine) receive() (r Reply, err error) {
	e.input.Reset()
	hdr := &header{}

	// decode header
	if err = hdr.read(e.reader); err != nil {
		if e.dumpConversation {
			fmt.Printf("%d< %v\n", e.client, err)
		}
		return
	}

	// decode message
	r, err = code2Msg(hdr.code)
	if err != nil {
		if e.dumpConversation {
			fmt.Printf("%d< %v %s\n", e.client, hdr, err)
		}
		return
	}

	if err = r.read(e.reader); err != nil {
		if e.dumpConversation {
			fmt.Printf("%d< %v %s\n", e.client, hdr, err)
		}
		return
	}

	if e.dumpConversation {
		dump := hdr.code != e.lastDumpRead
		e.lastDumpRead = hdr.code

		dump = dump || r.code() == mErrorMessage

		if mr, ok := r.(MatchedReply); ok {
			dump = dump || mr.Id() != e.lastDumpId
			e.lastDumpId = mr.Id()
		}

		if dump {
			str := fmt.Sprintf("%v", r)
			cut := len(str)
			if cut > 80 {
				str = str[:76] + "..."
			}
			fmt.Printf("%d< %v %s\n", e.client, hdr, str)
		}
	}

	return
}

type EngineState int

const (
	EngineReady EngineState = 1 << iota
	EngineExitError
	EngineExitNormal
)

func (s EngineState) String() string {
	switch s {
	case EngineReady:
		return "EngineReady"
	case EngineExitError:
		return "EngineExitError"
	case EngineExitNormal:
		return "EngineExitNormal"
	}
	panic("unreachable")
}
