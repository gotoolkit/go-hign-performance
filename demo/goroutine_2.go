package main

import (
	"time"
	"log"
)

func main() {
	c := make(chan int)
	for {
		time.Sleep(time.Second)
		go func() {
			log.Println("waiting")
			<-c
			log.Println("done")
		}()


		log.Println("loop")
	}

}