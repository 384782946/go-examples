package hello

// gomobile bind -target=android github.com/384782946/go-examples/test-gomobile

// gomobile bind -target=ios github.com/384782946/go-examples/test-gomobile

import (
	"fmt"
)

func SayHello(name string) error {
	fmt.Println("hello,", name)
	return nil
}
