package chanpiper_test

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transcelestial/chanpiper"
)

func TestChanpiper(t *testing.T) {
	defer leaktest.Check(t)

	source := make(chan chanpiper.Data)
	piper := chanpiper.New(source)

	var wg sync.WaitGroup
	pipeCount := 10
	dataBufSize := 20

	var mux sync.Once
	data := make(chan chanpiper.Data, dataBufSize)
	var cur int32
	done := make(chan struct{})

	for i := 0; i < pipeCount; i++ {
		wg.Add(1)
		c := piper.Pipe()
		go func() {
			defer wg.Done()
			for d := range c {
				data <- d
				atomic.AddInt32(&cur, 1)
				v := atomic.LoadInt32(&cur)
				fmt.Println(v)
				if v == int32(dataBufSize) {
					mux.Do(func() {
						close(done)
					})
				}
			}
		}()
	}

	runtime.Gosched()                  // doesn't work as we expect, hence the
	time.Sleep(500 * time.Microsecond) // enough time for the above routine to start receiving

	source <- chanpiper.Data{"test a"}
	source <- chanpiper.Data{"test b"}

	<-done
	close(source)

	wg.Wait()

	require.Len(t, data, dataBufSize)
	assert.Nil(t, piper.Pipe())
}
