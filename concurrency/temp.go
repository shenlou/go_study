package main

import (
	// "fmt"
	"log"
	"os"
	// "runtime"
	"time"
)

func Add(x, y int, c chan int) {

	c <- x + y
}

func main() {
	// log.Println(runtime.NumCPU())
	// runtime.GOMAXPROCS(runtime.NumCPU())
	logfile, err := os.OpenFile("temp.log", os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}
	defer logfile.Close()
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println(time.Now())
	max := 1000000
	ch := make(chan int, max)
	for i := 0; i < max; i++ {
		go Add(i, i, ch)
	}
	for i := 0; i < max; i++ {
		value, ok := <-ch
		if ok {
			close(ch)
		}
		logger.Println(value)
	}

	// for _, ch := range 50000 {
	// 	value := <-ch
	// 	fmt.Println(value)
	// }
	logger.Println(time.Now())
}
