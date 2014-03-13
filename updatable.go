package trade

import (
	"errors"
	"fmt"
	"time"
)

// Updatable provides an asynchronous, high-level abstraction over selected TWS
// API use cases. There are currently Updatables for metadata, options,
// instruments and option chains.
//
// An Updatable is constructed (via its NewXXX function) with request information
// relevant to the particular TWS API use case. No Engine calls are made at the
// time of construction.
//
// Invoking StartUpdate() will request the Updatable to communicate with the
// Engine.
//
// After Engine starts receiving Reply values from TWS, it will deliver these to
// the Updatable. The Updatable will consume those Reply values and signal the
// updates channel when a new piece of data is ready. The updates channel is
// available by calling Update().
//
// The caller can then invoke use case-specific accessor methods provided by the
// Updatable. These methods guarantee to return consistent views of the data, with
// consistency depending on the specific Updatable use case being abstracted.
//
// If any error occurs while receiving Reply values, the errors channel will be
// signalled. The state of an Updatable after signalling errors is unspecified.
// Callers should consider monitoring Engine state changes directly so they
// become aware if the Engine exits. In such a case there is no guarantee an
// Updatable will directly discover this (eg if the Updatable has elected not to
// subscribe to Engine state changes itself).
//
// Callers can stop further updates using StopUpdate(). Most Updatables also
// provide a Cleanup() method for destroy actions.
type Updatable interface {
	StartUpdate() error
	StopUpdate()
	Update() chan bool
	Error() chan error
}

func WaitForUpdate(v Updatable, timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		return updateError(v, errors.New("Update timeout"))
	case err := <-v.Error():
		return updateError(v, err)
	case <-v.Update():
	}
	return nil
}

type UpdateError struct {
	v   Updatable
	err error
}

func (u *UpdateError) Error() string {
	return fmt.Sprintf("Error %s while updating %v", u.err, u.v)
}

func updateError(v Updatable, err error) error {
	return &UpdateError{v, err}
}

type Updatables struct {
	items  []Updatable
	update chan bool
	error  chan error
}

func NewUpdatables() *Updatables {
	return &Updatables{
		items:  make([]Updatable, 0),
		update: make(chan bool),
		error:  make(chan error),
	}
}

func (u *Updatables) Update() chan bool { return u.update }
func (u *Updatables) Error() chan error { return u.error }

func (u *Updatables) Add(v Updatable) {
	u.items = append(u.items, v)
}

func (u *Updatables) StartUpdate(timeout time.Duration) error {
	for _, v := range u.items {
		if err := v.StartUpdate(); err != nil {
			return updateError(v, err)
		}
	}

	go func() {
		for _, v := range u.items {
			select {
			case <-time.After(timeout):
				u.error <- updateError(v, errors.New("Update timeout"))
				return
			case err := <-v.Error():
				u.error <- updateError(v, err)
				return
			case <-v.Update():
			}
		}

		u.update <- true
	}()

	return nil
}

func (u *Updatables) StopUpdate() {
	for _, v := range u.items {
		v.StopUpdate()
	}
}
