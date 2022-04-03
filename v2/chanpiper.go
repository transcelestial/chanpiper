// Package for piping data from a source chan to many
// receivers.
//
// Usage:
//  source := make(chan string)
//  piper := chanpiper.New(source)
//
//  go func(){
//    for data := range piper.Pipe() {
//      fmt.Println(data)
//    }
//    fmt.Println("done")
//  }()
//
//  go func(){
//    for data := range piper.Pipe() {
//      fmt.Println(data)
//    }
//    fmt.Println("done")
//  }()
//
//  // wait a bit before sending
//  time.Sleep(time.Second)
//
//  source <- "ping"
//  close(source)
package chanpiper

import (
	"sync"
)

// New returns a new chan piper interface.
func New[T any](source <-chan T) Piper[T] {
	p := &piper[T]{
		chans: []chan<- T{},
	}
	go p.setupRcvLoop(source)
	return p
}

type Piper[T any] interface {
	// Pipe returns a new chan that receives whenever the source sends data.
	//
	// The backpressure mechanism allows for bufferring 1 msg; so if the source sends too fast,
	// any extra msgs are discarded.
	//
	// Also note that there's no replay mechanism,
	// so any msgs that were sent before calling Pipe() are not received.
  //
  // Returns nil if the source is closed.
	Pipe() <-chan T
}

type piper[T any] struct {
	mux   sync.RWMutex
	chans []chan<- T
}

func (p *piper[T]) Pipe() <-chan T {
	p.mux.Lock()
	defer p.mux.Unlock()
	c := make(chan T, 1)
  if p.chans == nil {
		return nil
	}
	p.chans = append(p.chans, c)
	return c
}

func (p *piper[T]) setupRcvLoop(source <-chan T) {
	for d := range source {
		p.mux.RLock()
		for _, c := range p.chans {
			select {
			// don't block
			case c <- d:
			default:
			}
		}
		p.mux.RUnlock()
	}
  p.mux.Lock()
	defer p.mux.Unlock()
	// source got closed, so close the receivers
	for _, c := range p.chans {
		close(c)
	}
  // empty chans
	p.chans = nil
}
