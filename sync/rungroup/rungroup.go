package rungroup

import "sync"

// Go increments sync.WaitGroup, runs fn in a new goroutine and
// decrements sync.WaitGroup after fn returns.
func Go(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() { defer wg.Done(); fn() }()
}

// RunGroup is a sync.WaitGroup with a Go method.
type RunGroup struct {
	sync.WaitGroup
}

// Go increments RunGroup, runs fn in a new goroutine and
// decrements RunGroup after fn returns.
func (g *RunGroup) Go(fn func()) { Go(&g.WaitGroup, fn) }
