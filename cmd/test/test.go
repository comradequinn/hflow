package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		time.After(time.Millisecond * 500)
		fmt.Printf("\r")
	}

}
