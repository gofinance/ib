package trade

import (
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
		return updateError(v, timeoutError())
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

func (self *UpdateError) Error() string {
	return fmt.Sprintf("Error %s while updating %v", self.err, self.v)
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

func (self *Updatables) Update() chan bool { return self.update }
func (self *Updatables) Error() chan error { return self.error }

func (self *Updatables) Add(v Updatable) {
	self.items = append(self.items, v)
}

func (self *Updatables) StartUpdate(timeout time.Duration) error {
	for _, v := range self.items {
		if err := v.StartUpdate(); err != nil {
			return updateError(v, err)
		}
	}

	go func() {
		for _, v := range self.items {
			select {
			case <-time.After(timeout):
				self.error <- updateError(v, timeoutError())
				return
			case err := <-v.Error():
				self.error <- updateError(v, err)
				return
			case <-v.Update():
			}
		}

		self.update <- true
	}()

	return nil
}

func (self *Updatables) StopUpdate() {
	for _, v := range self.items {
		v.StopUpdate()
	}
}
