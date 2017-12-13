package main

import (
	"fmt"
	"time"
)

func timer(d time.Duration) <-chan int {
	c := make(chan int)
	go func() {
		time.Sleep(d)
		c <- 1
	}()
	return c
}

func main() {

	var ch chan int = make(chan int)
	go func(ch chan int){
		for i:=0; i<10; i++ {
			ch <- i*10
		}
		timer(13*time.Second)
	}(ch)
	for i := 0; i < 11; i++ {
		c := timer(1 * time.Second)
		<-c
		fmt.Println(i)
		fmt.Println(<-ch)
	}
}
