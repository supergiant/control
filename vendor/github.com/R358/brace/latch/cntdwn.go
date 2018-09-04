// Copyright (c) 2014 R358.org http://www.R358.org
package latch

import (
	"time"
	"sync"
	"sync/atomic"
)

//
// Implements a count down latch with timeouts.
// Any thread calling Await[OrTimeout|Func] will block until the latch is canceled, count <=0 or timeout occurs.
// Note: 1. TimedOut(), Canceled(), Completed() are only guaranteed to be consistent with each other after an await
//          method has finished.
//		 2. CountDown() and Cancel() are synchronized, Canceled() and Completed() will only be set after and await
// 			function has been called.
//       3. Only one thread may call any await method it is synchronized, second and subsequent will cause panic
// 			immediately or after blocking.
// 		 4. Do not reuse.
//
//
type CountdownLatch struct {
	count                 int32
	zeroReached           chan bool
	cancel                chan bool
	canceled              int32
	timedOut              int32
	completed             int32
	awaitMutex *sync.Mutex
	countMutex *sync.Mutex
	onComplete func(*CountdownLatch, bool)
}


type LatchError struct {
	LatchCompleted   bool
	CountReachedZero bool
	LatchCanceled    bool
}

func (l LatchError) Error() string {
	switch {
	case l.CountReachedZero:
		return "Count reached zero"
	case l.LatchCompleted:
		return "Latch has completed."

	case l.LatchCanceled:
		return "Latch has been canceled."

	default:
		return "Latch unknown error."
	}
}

func NewCountdownLatch(count int) *CountdownLatch {
	return &CountdownLatch{ int32(count), make(chan bool, 1), make(chan bool, 1), 0, 0, 0, new(sync.Mutex), new(sync.Mutex), nil}
}


func NewCountdownLatchWithListener(count int, onComplete func(*CountdownLatch, bool)) *CountdownLatch {
	return &CountdownLatch{ int32(count), make(chan bool, 1), make(chan bool, 1), 0, 0, 0, new(sync.Mutex), new(sync.Mutex), onComplete}
}

//
// TimedOut, not guaranteed to be consistent until after Await has completed.
//
func (c *CountdownLatch) TimedOut() bool {
	return atomic.LoadInt32(&c.timedOut) == 1
}

//
// Canceled, not guaranteed to be consistent until after Await has completed.
//
func (c *CountdownLatch) Canceled() bool {
	return atomic.LoadInt32(&c.canceled) == 1
}

//
// Call to decrement the latch.
//
func (c *CountdownLatch) CountDown() {
	c.countMutex.Lock()
	defer c.countMutex.Unlock()

	if c.Completed() {
		panic(LatchError{LatchCompleted:true})
	}

	if c.Canceled() {
		panic(LatchError{LatchCanceled: true})
	}

	if c.count == 0 {
		panic(LatchError{CountReachedZero:true})
	}

	c.count--
	if c.count <= 0 {
		c.zeroReached <- true
		if c.onComplete != nil {
			c.onComplete(c, false)
			c.onComplete = nil // This protected by c.countMutex
		}
		c.count = 0
	}

}

//
// Call to cancel the latch, cancel can come from any thread.
//
func (c *CountdownLatch) Cancel() {
	c.countMutex.Lock()
	defer c.countMutex.Unlock()
	atomic.StoreInt32(&c.canceled, 1)
	c.cancel <- true
	if c.onComplete != nil {
		c.onComplete(c, true)
		c.onComplete = nil // This protected by c.countMutex
	}
}

//
// Call to to determine if latch reached zero, consistent only after Await[Timeout|Func] has completed.
//
func (c *CountdownLatch) Completed() bool {
	return atomic.LoadInt32(&c.completed) == 1
}

//
// Await countdown to <=0 or return true if canceled
//
func (c *CountdownLatch) Await() bool {

	c.awaitMutex.Lock()
	defer func() {
		atomic.StoreInt32(&c.completed, 1)
		c.awaitMutex.Unlock();
	}()


	if c.Completed() {
		panic(LatchError{LatchCompleted:true})
	}

	for {
		select {

		case <-c.cancel:
			//			atomic.StoreInt32(&c.canceled, 1)
			return true


		case <-c.zeroReached:
			atomic.StoreInt32(&c.completed, 1)
			return false
			break
		}
	}
}

//
// AwaitOrFail, the count must be <= 0 before the duration expires.
// Will Return true on timeout or cancel.
//
func (c *CountdownLatch) AwaitTimeout(duration time.Duration) bool {

	c.awaitMutex.Lock()
	defer func() {
		atomic.StoreInt32(&c.completed, 1)
		c.awaitMutex.Unlock();
	}()

	if c.Completed() {
		panic(LatchError{LatchCompleted:true})
	}

	to := time.After(duration)


	for {
		select {

		case <-to:
			atomic.StoreInt32(&c.timedOut, 1)
			return true

		case <-c.cancel:
			atomic.StoreInt32(&c.canceled, 1)
			return true

		case <-c.zeroReached:
			atomic.StoreInt32(&c.completed, 1)
			return false
			break
		}
	}
}

//
// AwaitFunc will call onZero when count <=0, will call onTimeout if the timeout has expired and onCancel if it is canceled.
//
func (c *CountdownLatch) AwaitFunc(duration time.Duration, onZero func(*CountdownLatch), onCancel func(*CountdownLatch), onTimeout func(*CountdownLatch)) bool {
	c.awaitMutex.Lock()
	defer c.awaitMutex.Unlock()

	to := time.After(duration)

	if c.Completed() {
		panic(LatchError{LatchCompleted:true})
	}


	for {
		select {
		case <-to:
			atomic.StoreInt32(&c.timedOut, 1)
			if onTimeout != nil {
				onTimeout(c)
				return true
			}

		case <-c.cancel:
			atomic.StoreInt32(&c.canceled, 1)
			if onCancel != nil {
				onCancel(c)
				return true
			}

		case <-c.zeroReached:
			atomic.StoreInt32(&c.completed, 1)
			if onZero != nil {
				onZero(c)
			}
			return false
			break
		}
	}
}
