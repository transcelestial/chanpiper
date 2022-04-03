// Package for piping data from a source chan to many
// receivers.
//
// Usage:
//  source := make(chan chanpiper.Data)
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
//  source <- chanpiper.Data{"ping"}
//  close(source)
package chanpiper

import (
	"sync"
)

// New creates a new chan piper.
func New(source <-chan Data) Piper {
	p := &piper{
		chans: []chan<- Data{},
	}
	go p.setupRcvLoop(source)
	return p
}

type Piper interface {
	// Pipe returns a new chan that receives whenever the source sends data.
	//
	// The backpressure mechanism allows for bufferring 1 msg; so if the source sends too fast,
	// any extra msgs are discarded.
	//
	// Also note that there's no replay mechanism,
	// so any msgs that were sent before calling Pipe(), are not received.
	Pipe() <-chan Data
}

type Data struct {
	V interface{}
}

type piper struct {
	mux   sync.RWMutex
	chans []chan<- Data
}

func (p *piper) Pipe() <-chan Data {
	p.mux.Lock()
	defer p.mux.Unlock()
	c := make(chan Data, 1)
	p.chans = append(p.chans, c)
	return c
}

func (p *piper) setupRcvLoop(source <-chan Data) {
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
	// source got closed, so close the receivers
	for _, c := range p.chans {
		close(c)
	}
}
