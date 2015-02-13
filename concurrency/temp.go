package main

import (
	// "fmt"
	"log"
	"runtime"
)

func Add(c chan int) {
	a := 0
	for i := 0; i < 100; i++ {
		a += i
		log.Printf("print a %d", a)
	}

	c <- 1
}

func Print(c chan int) {
	for i := 0; i < 100; i++ {
		log.Printf("print d %d", i)
	}
	c <- 1
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan int, 2)
	go Print(ch)
	go Add(ch)

	for i := 0; i < 2; i++ {
		<-ch
	}

}
