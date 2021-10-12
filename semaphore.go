package qps_limiter
// Concurrency Control
type Semaphore struct {
	ch chan struct{}
}

// Create semaphore
func New(size int) *Semaphore {
	if size < 1 {
		size = 1
	}
	return &Semaphore{ch: make(chan struct{}, size)}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) Size() int {
	return cap(s.ch)
}

func (s *Semaphore) Free() int {
	return cap(s.ch) - len(s.ch)
}