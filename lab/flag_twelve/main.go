package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("flag{12}")
			os.Exit(0)
		}
	}()

	for {
		time.Sleep(3 * time.Second)
	}
}
