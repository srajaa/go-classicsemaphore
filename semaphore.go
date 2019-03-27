package csemaphore // import "github.com/srajaa/go-classical-semaphore/csemaphore"

import "sync"

//CSemaphore can be zero initialized and
//should be ready to go, or it can be initialized
//with a positive count
type CSemaphore struct {
	count   uint64
	mutex   sync.Mutex
	waiters [](chan int)
}

//P blocks the caller if semaphore count is zero,
//or just goes through
func (s *CSemaphore) P() {
	s.mutex.Lock()
	if s.count != 0 && len(s.waiters) > 0 {
		panic("Unexpected state")
	}
	if s.count < 0 {
		panic("Unexpected state")
	}
	if s.count > 0 {
		s.count--
		s.mutex.Unlock()
		return
	}
	waiter := make(chan int, 1)
	s.waiters = append(s.waiters, waiter)
	s.mutex.Unlock()
	<-waiter
}

//V releases one of the waiters in FIFO order or
//or just increments the count
func (s *CSemaphore) V() {
	s.mutex.Lock()
	if s.count != 0 && len(s.waiters) > 0 {
		panic("Unexpected state")
	}
	if s.count < 0 {
		panic("Unexpected state")
	}
	if len(s.waiters) > 0 {
		waiter := s.waiters[0]
		s.waiters = s.waiters[1:]
		s.mutex.Unlock()
		waiter <- 1
		return
	}
	s.count++
	s.mutex.Unlock()
}
