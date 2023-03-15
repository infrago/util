package util

import (
	"sync"
)

type Runner struct {
	cs chan error
	wg sync.WaitGroup
}

func NewRunner() *Runner {
	s := &Runner{
		cs: make(chan error),
	}

	return s
}

func (s *Runner) Run(f func()) {
	s.wg.Add(1)
	go func() {
		f()
		s.wg.Done()
	}()
}

func (s *Runner) Stop() chan error {
	return s.cs
}

func (s *Runner) Wait() {
	s.wg.Wait()
}
func (s *Runner) End() {
	close(s.cs)
	s.wg.Wait()
}
