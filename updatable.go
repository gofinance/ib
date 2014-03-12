package trade

import (
	"errors"
	"fmt"
	"time"
)

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
