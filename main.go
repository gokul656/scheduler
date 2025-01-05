package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	workerpool := NewWorkerPool(5)
	workerpool.Spawn()

	for range 100 {
		workerpool.Submit(doSomething)
	}

	workerpool.Stop()
}

func doSomething() error {
	time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	log.Println("doing something")
	return nil
}
