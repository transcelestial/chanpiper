package chanpiper_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/transcelestial/chanpiper"
)

func Example() {
	source := make(chan chanpiper.Data)
	piper := chanpiper.New(source)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range piper.Pipe() {
			fmt.Println(data.V)
		}
		fmt.Println("done")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range piper.Pipe() {
			fmt.Println(data.V)
		}
		fmt.Println("done")
	}()

	// wait a bit before sending
	time.Sleep(100 * time.Millisecond)

	// send some data
	source <- chanpiper.Data{"ping"}

	// wait a bit more before we close
	time.Sleep(100 * time.Millisecond)

	// close
	close(source)

	wg.Wait()

	// Output:
	// ping
	// ping
	// done
	// done
}
