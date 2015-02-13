package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Add(x, y int) {

	z := x + y
	fmt.Println(z)
}

func main() {
	logfile, err := os.OpenFile("temp.log", os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}
	defer logfile.Close()
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println(time.Now())
	for i := 0; i < 500000; i++ {
		Add(i, i)
	}
	logger.Println(time.Now())
}
