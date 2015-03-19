package rtda

import (
	"time"
)

const (
	_wakeUp    = 1
	_interrupt = 2
)

func (self *Thread) Sleep(d time.Duration) (interrupted bool) {
	self.lock.Lock()
	if self.isInterrupted {
		self.isInterrupted = false
		self.lock.Unlock()
		return true
	}

	self.isBlocked = true
	go self._sleep(d)
	self.lock.Unlock()

	interrupted = (<-self.ch == _interrupt)
	return
}

func (self *Thread) _sleep(d time.Duration) {
	time.Sleep(d)

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.isBlocked { // not interrupted
		self.isBlocked = false
		self.ch <- _wakeUp
	}
}
