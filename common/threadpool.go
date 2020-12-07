package resources

import (
  "sync/atomic"
  "time"
)

// ThreadPool manages the active thread pool
type ThreadPool struct {
  maxCount     int32
  activeCount  int32
  minWaitTimer time.Duration
}

// NewThreadPool : creates a new threadpool
func NewThreadPool(maxThreadCount int32) (pool ThreadPool) {
  pool = ThreadPool{maxThreadCount, int32(0), 1 * time.Millisecond}
  return
}

// SetWaitTimer : change actual wait timer (1 Millisecond)
func (pg *ThreadPool) SetWaitTimer(multiplier time.Duration) {
  pg.minWaitTimer = multiplier
}

// Add worker
func (pg *ThreadPool) Add(n int32) {
  atomic.AddInt32(&pg.activeCount, n);
}

// Done :: Adds -1
func (pg *ThreadPool) Done() {
  pg.Add(-1)
}

// Wait until the activeCount is 0
func (pg *ThreadPool) Wait() {
  for ;atomic.LoadInt32(&pg.activeCount) > 0; {
    time.Sleep(pg.minWaitTimer)
  }
}

// WaitWithTimeout : wait with a max timeout. 1 Tick = total time per each minWaitTimer
// returns true if wait returns without timeout, otherwise false
func (pg *ThreadPool) WaitWithTimeout(tickCount int) bool {
  count := 0
  for ;atomic.LoadInt32(&pg.activeCount) > 0; {
    time.Sleep(pg.minWaitTimer)
    count++
    if count >= tickCount { return false }
  }
  return true
}

// WaitPool : waits until the pool is available
func (pg *ThreadPool) WaitPool() {
  for ;atomic.LoadInt32(&pg.activeCount) >= pg.maxCount; {
    time.Sleep(pg.minWaitTimer)
  }
}
