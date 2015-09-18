package ib

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// UpdateStatus .
type UpdateStatus int

// Status enum
const (
	UpdateFalse UpdateStatus = 1 << iota
	UpdateTrue
	UpdateFinish
)

// Manager provides a high-level abstraction over selected IB API use cases.
// It defines a contract for executing IB API operations with explicit
// isolation, error handling and concurrency guarantees.
//
// To use a Manager, invoke its NewXXManager function. This will immediately
// return an XXManager value that will subsequentially asynchronously connect to
// the Engine and send data. The only errors returned from a NewXXManager
// function relate to invalid parameters, as Engine state is unknown until the
// asynchronous connection is attempted at an unspecified later time.
//
// A Manager provides idempotent accessors to obtain data or the state of the
// construction-specified operation. Accessors guarantee to never return data
// from a partially-consumed IB API Reply. They may return zero or nil if data
// isn't yet available. Accessors remain usable even after a Manager is closed.
//
// Clients must invoke the Refresh() method to obtain the channel that will be
// signalled whenever the Manager has new data available from one or more
// accessors. The Manager will close the refresh channel when no further updates
// will be provided. This may occur for three reasons: (i) a client has invoked
// Close(); (ii) the Manager encountered an error (which is available via
// FatalError()) or (iii) the Manager is aware no new data will be sent by this
// Manager (typically as the IB API has sent an "end of data" Reply). In either
// of the latter two situations the Manager will automatically unsubscribe from
// IB API subscriptions and Engine subscriptions (this is a performance
// optimisation only; users should still use "defer manager.Close()").
//
// A Manager will block until a client can consume from the refresh channel.
// If the client needs to perform long-running operations, it should consider
// dedicating a goroutine to consuming from the channel and maintaining a count
// of update signals. This would permit the client to discover it missed an
// update without blocking the Manager (or its upstream Engine). SinkManager is
// another option for clients only interested in the final result of a Manager.
//
// Managers will register for Engine state changes. If the Engine exits, the
// Manager will close. If Engine.FatalError() returns an error, it will be made
// available via manager.FatalError() (unless an earlier error was recorded).
// This means clients need not track Manager errors or states themselves.
//
// Every Manager defines a Close() method which blocks until the Manager has
// released its resources. The Close() method will not return any new error or
// change the state of FatalError(), even if it encounters errors (eg Engine
// send failure) while closing its resources. After a Manager is closed, its
// accessors remain available. Repeatedly calling Close() is safe and will not
// repeatedly release any already-released resources. Therefore it is safe and
// highly recommended to use "defer manager.Close()".
type Manager interface {
	FatalError() error
	Refresh() <-chan bool
	Close()
}

// AbstractManager implements most of the Manager interface contract.
type AbstractManager struct {
	rwm    sync.RWMutex
	term   chan struct{}
	exit   chan bool
	update chan bool
	engs   chan EngineState
	eng    *Engine
	err    error
	rc     chan Reply
}

// NewAbstractManager .
func NewAbstractManager(e *Engine) (*AbstractManager, error) {
	if e == nil {
		return nil, errors.New("Engine required")
	}
	return &AbstractManager{
		rwm:    sync.RWMutex{},
		term:   make(chan struct{}),
		exit:   make(chan bool),
		update: make(chan bool),
		engs:   make(chan EngineState),
		eng:    e,
		rc:     make(chan Reply),
	}, nil
}

func (a *AbstractManager) RecvChan() chan Reply {
	return a.rc
}

func (a *AbstractManager) Engine() *Engine {
	return a.eng
}

func (a *AbstractManager) StartMainLoop(preLoop func() error, receive func(r Reply) (UpdateStatus, error), preDestroy func()) {
	defer func() {
		preDestroy()
		a.eng.UnsubscribeState(a.engs)
		close(a.update)
		close(a.term)
	}()

	go a.eng.SubscribeState(a.engs)
	if err := preLoop(); err != nil {
		a.err = err
		return
	}

	for {
		select {
		case <-a.exit:
			return
		case r := <-a.rc:
			if a.consume(r, receive) {
				return
			}
		case <-a.engs:
			if a.err == nil {
				a.err = a.eng.FatalError()
			}
			return
		}
	}
}

// consume handles sending one Reply to the receive function. Returning true
// indicates the main loop should terminate (ie the AbstractManager close).
func (a *AbstractManager) consume(r Reply, receive func(r Reply) (UpdateStatus, error)) (exit bool) {
	updStatus := make(chan UpdateStatus)

	go func() { // new goroutine to guarantee unlock
		a.rwm.Lock()
		defer a.rwm.Unlock()
		status, err := receive(r)
		if err != nil {
			a.err = err
			close(updStatus)
			return
		}
		updStatus <- status
	}()

	status, ok := <-updStatus
	if !ok {
		return true // channel closed due to receive func error result
	}
	switch status {
	case UpdateFalse:
	case UpdateTrue:
		a.update <- true
	case UpdateFinish:
		a.update <- true
		return true
	}
	return false
}

// FatalError .
func (a *AbstractManager) FatalError() error {
	return a.err
}

// Refresh .
func (a *AbstractManager) Refresh() <-chan bool {
	return a.update
}

// Close .
func (a *AbstractManager) Close() {
	select {
	case <-a.term:
		return
	case a.exit <- true:
	}
	<-a.term
}

// SinkManager listens to a Manager until it closes the update channel or reaches
// a target update count or timeout. This function is useful for clients who only
// require the final result of a Manager (and have no interest in each update).
// The Manager is guaranteed to be closed before it returns.
func SinkManager(m Manager, timeout time.Duration, updateStop int) (updates int, err error) {
	for {
		sentClose := false
		select {
		case <-time.After(timeout):
			m.Close()
			return updates, fmt.Errorf("SinkManager: no new update in %s", timeout)
		case _, ok := <-m.Refresh():
			if !ok {
				return updates, m.FatalError()
			}
			updates++
			if updates >= updateStop && !sentClose {
				sentClose = true
				go m.Close()
			}
		}
	}
}
