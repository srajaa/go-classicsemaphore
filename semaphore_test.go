package csemaphore

import (
	"log"
	"sync"
	"testing"
)

type itemsstruct struct {
	mutex sync.Mutex
	count int64
}

func producer(it *itemsstruct, s *CSemaphore, barrierCh chan int) {
	it.mutex.Lock()
	it.count++
	if it.count > 100 {
		panic("bad count")
	}
	log.Println("Produced", it.count)
	it.mutex.Unlock()
	s.V()
	barrierCh <- 1
}

func consumer(it *itemsstruct, s *CSemaphore, barrierCh chan int) {
	s.P()
	it.mutex.Lock()
	it.count--
	if it.count < 0 {
		panic("bad count")
	}
	log.Println("Consumed", it.count)
	it.mutex.Unlock()
	barrierCh <- 1
}

func Test_ProdConsumer(t *testing.T) {
	barrierCh := make(chan int, 200)
	it := itemsstruct{}
	s := CSemaphore{}
	for i := 0; i < 100; i++ {
		go producer(&it, &s, barrierCh)
	}
	for i := 0; i < 100; i++ {
		go consumer(&it, &s, barrierCh)
	}
	for i := 0; i < 200; i++ {
		<-barrierCh
	}
	if it.count != 0 {
		t.Error("Semaphore test failed")
	}
	for i := 0; i < 100; i++ {
		go consumer(&it, &s, barrierCh)
	}
	for i := 0; i < 100; i++ {
		go producer(&it, &s, barrierCh)
	}
	for i := 0; i < 200; i++ {
		<-barrierCh
	}
	if it.count != 0 {
		t.Error("Semaphore test failed")
	}
}
