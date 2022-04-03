# chanpiper
> A simple interface that implements *-to-many data flow using channels in Go.

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/transcelestial/chanpiper/Test?label=test&style=flat-square)](https://github.com/transcelestial/chanpiper/actions?query=workflow%3ATest)

## Usage
To use `chanpiper` w/o generics:
```go
package main

import (
    "fmt"

    "github.com/transcelestial/chanpiper"
)

func main() {
    source := make(chan chanpiper.Data)
	piper := chanpiper.New(source)

    p1 := piper.Pipe()
    p2 := piper.Pipe()

    go func() {
		for data := range p1 {
			fmt.Println(data.V)
		}
	}()

    go func() {
		for data := range p2 {
			fmt.Println(data.V)
		}
	}()

    source <- chanpiper.Data{"ping"}
    close(source)
}
```
See the [example](./example_chanpiper_test.go) for a working example.

To use `chanpiper/v2` w/ generics:
```go
package main

import (
    "fmt"

    "github.com/transcelestial/chanpiper/v2"
)

func main() {
    source := make(chan string)
	piper := chanpiper.New(source)

    p1 := piper.Pipe()
    p2 := piper.Pipe()

    go func() {
		for data := range p1 {
			fmt.Println(data.V)
		}
	}()

    go func() {
		for data := range p2 {
			fmt.Println(data.V)
		}
	}()

    source <- "ping"
    close(source)
}
```
See the [example](./v2/example_chanpiper_test.go) for a working example.

## Contribute
If you wish to contribute, please use the following guidelines:
* Use [conventional commits](https://conventionalcommits.org/)
* Use [effective Go](https://golang.org/doc/effective_go)
