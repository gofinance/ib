// Package ib trade offers a pure Go abstraction over Interactive Brokers IB API.
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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"time"
)

// consts .
const (
	gatewayDefault   = "127.0.0.1:4001"
	UnmatchedReplyID = int64(-9223372036854775808)
	maxMsgLength     = 0xffffff
)

// EngineOptions .
type EngineOptions struct {
	Gateway          string
	Client           int64
	DumpConversation bool
	DisableV100Plus  bool
	Logger           *log.Logger
}

// Engine is the entry point to the IB TWS API
type Engine struct {
	id               chan int64
	exit             chan bool
	terminated       chan struct{}
	ch               chan command
	gateway          string
	client           int64
	connected        bool
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
	lastDumpID       int64
	fatalError       error
	useV100Plus      bool
	logger           *log.Logger
}

type command struct {
	fun func()
	ack chan struct{}
}

type txrequest struct {
	req Request
	ack chan struct{}
}

func uniqueID(start int64) chan int64 {
	ch := make(chan int64)
	id := start
	go func() {
		for {
			if id == UnmatchedReplyID {
				id++
			}
			ch <- id
			id++
		}
	}()
	return ch
}

// Next client id. Package scope to ensure new engines have unique client IDs.
var clientSeq = uniqueID(1)

// NewEngine takes a client id and returns a new connection
// to IB Gateway or IB Trader Workstation.
func NewEngine(opt EngineOptions) (re *Engine, rerr error) {
	gateway := opt.Gateway
	if gateway == "" {
		gateway = gatewayDefault
	}

	logger := opt.Logger
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	conn, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rerr != nil {
			_ = conn.Close()
		}
	}()

	client := opt.Client

	e := Engine{
		id:               uniqueID(100),
		exit:             make(chan bool),
		terminated:       make(chan struct{}),
		ch:               make(chan command),
		gateway:          gateway,
		client:           client,
		con:              conn,
		reader:           bufio.NewReader(conn),
		input:            bytes.NewBuffer(make([]byte, 0, 4096)),
		output:           bytes.NewBuffer(make([]byte, 0, 4096)),
		rxReply:          make(chan Reply),
		rxErr:            make(chan error),
		txRequest:        make(chan txrequest),
		txErr:            make(chan error),
		observers:        map[int64]chan<- Reply{},
		state:            EngineReady,
		dumpConversation: opt.DumpConversation,
		useV100Plus:      !opt.DisableV100Plus,
		logger:           logger,
	}

	if err := e.handshake(); err != nil {
		return nil, err
	}

	// start worker goroutines (these exit on request or error)
	go e.startReceiver()
	go e.startTransmitter()
	go e.startMainLoop()

	e.logger.Printf("Client: %v", e.client)
	// send the StartAPI request
	e.Send(&StartAPI{Client: e.client})

	e.connected = true

	return &e, nil
}

// IsUseV100Plus .
func (e *Engine) IsUseV100Plus() bool {
	return e.useV100Plus
}

func (e *Engine) makeV100APIHeader() {

	out := buildVersionString(mMinVersion, mMaxVersion)

	//if e.connectOptions != "" {
	//	out += " " + connectOptions
	//}

	e.output.Reset()
	e.output.WriteString("API\000")
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(len(out)))
	e.output.Write(data)
	e.output.WriteString(out)
}

func buildVersionString(minVersion, maxVersion int) string {
	if minVersion < maxVersion {
		return fmt.Sprintf("v%v..%v", minVersion, maxVersion)
	}

	return fmt.Sprintf("v%v", minVersion)
}

func (e *Engine) handshake() error {
	e.output = bytes.NewBuffer(make([]byte, 0, 4096))

	// send client version (unless logon via iserver and/or Version > 100)
	if !e.useV100Plus {
		// write client version
		// Do not add length prefix here, because Server does not know Client's version yet
		clientShake := &clientHandshake{clientVersion}
		if err := clientShake.write(e.output); err != nil {
			return err
		}
		e.logger.Printf("OLD")
	} else {
		// Switch to GW API (Version 100+ requires length prefix)
		e.makeV100APIHeader()
		e.logger.Printf("100V")
	}

	e.logger.Printf("WRITE")
	if _, err := e.con.Write(e.output.Bytes()); err != nil {
		return err
	}

	e.logger.Printf("READ")
	// read server version and time
	serverShake := &serverHandshake{}
	e.input.Reset()
	if e.useV100Plus {
		num, err := readUInt32(e.reader)
		if err != nil {
			return err
		}

		data := make([]byte, num)

		// TODO: check read size
		if _, err := e.reader.Read(data); err != nil {
			return err
		}

		if err = serverShake.read(e.serverVersion, bufio.NewReader(bytes.NewBuffer(data))); err != nil {
			return err
		}
	} else {
		if err := serverShake.read(e.serverVersion, e.reader); err != nil {
			return err
		}
	}

	if e.dumpConversation {
		e.logger.Printf("SERVER VERSION %v\n", serverShake.version)
	}

	if serverShake.version < minServerVersion {
		return fmt.Errorf("%s must be at least version %d (reported %d)", e.ConnectionInfo(), minServerVersion, serverShake.version)
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
			if err := e.transmit(t.req); err != nil {
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
					e.logger.Printf("Waited 5 seconds for state channel %v\n", ob)
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
			e.logger.Printf("%s engine: RX error %s", e.ConnectionInfo(), err)
			e.fatalError = err
			e.state = EngineExitError
			return
		case err := <-e.txErr:
			e.logger.Printf("%s engine: TX error %s", e.ConnectionInfo(), err)
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
	outer:
		for _, o := range e.observers {
			for _, prevDone := range done {
				if o == prevDone {
					continue outer
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
		if o, ok := e.observers[mr.ID()]; ok {
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
			e.logger.Printf("IB: Waited 5 seconds for reply channel %v, lost %v\n", c, r)
			return
		}
	}
}

func (e *Engine) transmit(r Request) (err error) {
	e.output = bytes.NewBuffer(make([]byte, 0, 4096))

	// encode the message itself
	if err = r.write(e.serverVersion, e.output); err != nil {
		return
	}

	if e.useV100Plus {
		data := make([]byte, 4)
		outputbytes := e.output.Bytes()
		if e.dumpConversation {
			e.logger.Printf("MESSAGE SIZE: %v\n", len(outputbytes))
		}

		binary.BigEndian.PutUint32(data, uint32(len(outputbytes)))
		outputbytes = append(data, outputbytes...)

		e.output = bytes.NewBuffer(outputbytes)
	}

	if e.dumpConversation {
		b := e.output
		s := strings.Replace(b.String(), "\000", "-", -1)
		e.logger.Printf("DUMP: %d> (%v) '%s'\n", e.client, e.output.Len(), s)
		e.logger.Printf("MESSAGE:\n%v\n", hex.Dump(b.Bytes()))
	}

	cnt, err := e.con.Write(e.output.Bytes())
	if e.dumpConversation {
		e.logger.Printf("DUMP: %d> wrote %v\n", e.client, cnt)
	}
	return
}

// NextRequestID returns a unique request id (which is never UnmatchedReplyID).
func (e *Engine) NextRequestID() int64 {
	return <-e.id
}

// ClientID .
func (e *Engine) ClientID() int64 {
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
		e.logger.Println("Engine unexpectedly terminated after command sent")
		return
	case <-cmd.ack:
		return
	}
}

// Subscribe will notify subscribers of future events with given id.
// Many request types implement MatchedRequest and therefore provide a SetID().
// To receive the corresponding MatchedReply events, firstly subscribe with the
// same id as will be assigned with SetID(). Any incoming events that do not
// implement MatchedReply will be delivered to those observers subscribed to
// the UnmatchedReplyID constant. Note that the engine will raise an error if
// an attempt is made to send a MatchedRequest with UnmatchedReplyID as its id,
// given the high unlikelihood of that id being required in normal situations
// and that NextRequestID() guarantees to never return UnmatchedReplyID.
// Each ErrorMessage event is delivered once only to each known observer.
// The engine never closes the channel (allowing reuse across IDs and engines).
// This call will block until the subscriber is registered or engine terminates.
func (e *Engine) Subscribe(o chan<- Reply, id int64) {
	e.sendCommand(func() {
		if id != UnmatchedReplyID {
			e.observers[id] = o
			return
		}

		e.unObservers = append(e.unObservers, o)
	})
}

// SubscribeAll .
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
		if id != UnmatchedReplyID {
			delete(e.observers, id)
			return
		}

		newUnObs := []chan<- Reply{}
		for _, existing := range e.unObservers {
			if existing != o {
				newUnObs = append(newUnObs, o)
			}
		}
		e.unObservers = newUnObs
	})
	close(terminate)
}

// UnsubscribeAll .
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
		newUnObs := []chan<- Reply{}
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

// Send a message to the engine, blocking until sent or the engine exits.
// This method will return an error if the UnmatchedReplyID is used or the
// engine exits. A nil error indicates successful transmission. Any transmission
// failure (eg connectivity loss) will cause the engine to exit with an error.
func (e *Engine) Send(r Request) error {
	if mr, ok := r.(MatchedRequest); ok {
		if mr.ID() == UnmatchedReplyID {
			return fmt.Errorf("%d is a reserved ID (try using NextRequestID)", UnmatchedReplyID)
		}
	}
	t := txrequest{r, make(chan struct{})}

	// send tx request
	select {
	case <-e.terminated:
		if err := e.FatalError(); err != nil {
			return err
		}
		return fmt.Errorf("Engine has already exited normally")
	case e.txRequest <- t:
	}

	// await ack or error
	select {
	case <-e.terminated:
		if err := e.FatalError(); err != nil {
			return err
		}
		return fmt.Errorf("Engine has already exited normally")
	case <-t.ack:
		return nil
	}
}

type packetError struct {
	value interface{}
	kind  reflect.Type
}

func (e *packetError) Error() string {
	return fmt.Sprintf("don't understand packet '%v' of type '%v'", e.value, e.kind)
}

func failPacket(v interface{}) error {
	return &packetError{
		value: v,
		kind:  reflect.ValueOf(v).Type(),
	}
}

func (e *Engine) receive() (Reply, error) {
	var reader *bufio.Reader

	if e.useV100Plus {
		msgSize, err := readUInt32(e.reader)
		if err != nil {
			return nil, err
		}

		if e.dumpConversation {
			e.logger.Printf("READ SIZE: %v\n", msgSize)
		}

		if msgSize > maxMsgLength {
			// TODO: handle this gracefully
			panic("message is too long: " + fmt.Sprintf("%v", msgSize))
		}

		data := make([]byte, msgSize)

		// TODO: check read size
		cnt, err := e.reader.Read(data)
		if err != nil {
			return nil, err
		}

		if e.dumpConversation {
			e.logger.Printf("READ SIZE1: %v of %v\n", cnt, msgSize)
		}

		// handle case where read does not fill the buffer
		for cnt != int(msgSize) {
			tmp, err := e.reader.Read(data[cnt:])
			if err != nil {
				return nil, err
			}
			cnt += tmp
			if e.dumpConversation {
				e.logger.Printf("READ SIZE2: %v of %v\n", cnt, msgSize)
			}
		}

		if e.dumpConversation {
			e.logger.Printf("READ MESSAGE:\n%v\n", hex.Dump(data))
		}

		e.input = bytes.NewBuffer(data)

		reader = bufio.NewReader(e.input)
	} else {
		e.input.Reset()
		reader = e.reader
	}

	code, err := readInt(reader)
	if err != nil {
		if e.dumpConversation {
			e.logger.Printf("DUMP: %d< %v\n", e.client, err)
		}
		return nil, err
	}

	// decode message
	r, err := code2Msg(code)
	if err != nil {
		if e.dumpConversation {
			e.logger.Printf("DUMP: %d< %v %s\n", e.client, code, err)
		}
		return nil, err
	}

	if err := r.read(e.serverVersion, reader); err != nil {
		if e.dumpConversation {
			e.logger.Printf("DUMP: %d< %v %s\n", e.client, code, err)
		}
		return nil, err
	}

	if e.dumpConversation {
		dump := code != e.lastDumpRead
		e.lastDumpRead = code

		dump = dump || r.code() == mErrorMessage

		if mr, ok := r.(MatchedReply); ok {
			dump = dump || mr.ID() != e.lastDumpID
			e.lastDumpID = mr.ID()
		}

		if dump {
			str := fmt.Sprintf("%v", r)
			cut := len(str)
			if cut > 80 {
				str = str[:76] + "..."
			}
			e.logger.Printf("DUMP: %d< %v %s\n", e.client, code, str)
		}
	}

	return r, nil
}

// EngineState .
type EngineState int

// Engine State enum
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
	default:
		panic("unreachable")
	}
}
