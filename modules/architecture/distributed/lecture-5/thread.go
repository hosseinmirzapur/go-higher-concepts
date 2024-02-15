package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	// simple concurrency
	concurrentDesign()

	// calling mulitple RPCs
	rpcDesign()

	// handling periodic functionality
	periodicDesign()

	// raft design which is calling multiple periodic functions
	raftDesign()

	// A simple bank transfer between alice and bob
	simpleBankTransfer()

	// conditional variable
	condVariable()

	// a simple usage of channels
	simpleChannel()

	// example of an RPC call via channel usage
	rpcWithChannels()

}

func concurrentDesign() {
	var a string

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		a = "Hello World"
		wg.Done()
	}()

	wg.Wait()
	fmt.Println(a)
}

func rpcDesign() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			fmt.Println(x)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func periodicDesign() {
	time.Sleep(1 * time.Second)

	fmt.Println("started")

	go periodic()

	time.Sleep(5 * time.Second)

}

func periodic() {
	for {
		fmt.Println("tick")
		time.Sleep(1 * time.Second)
	}
}

var done bool = false // initial state
var mu sync.Mutex

func raftDesign() {
	time.Sleep(1 * time.Second)
	fmt.Println("started")
	go periodicWithReturn()
	time.Sleep(5 * time.Second)

	mu.Lock()

	done = true

	mu.Unlock()
	fmt.Println("cancelled")

	time.Sleep(3 * time.Second)
}

func periodicWithReturn() {
	for {
		fmt.Println("tick")

		time.Sleep(1 * time.Second)
		mu.Lock()
		if done {
			return
		}

		mu.Unlock()
	}
}

func simpleBankTransfer() {
	alice := 10_000
	bob := 10_000

	total := alice + bob

	var mu sync.Mutex

	go func() {
		for i := 0; i < 1000; i++ {
			mu.Lock()
			defer mu.Unlock()
			alice -= 1
			bob += 1
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			mu.Lock()
			defer mu.Unlock()
			bob -= 1
			alice += 1
		}
	}()

	start := time.Now()

	for time.Since(start) < 1*time.Second {
		mu.Lock()
		defer mu.Unlock()
		if alice+bob != total {
			log.Printf("observed violation in bank transfer: alice -> %d  /  bob -> %d \n", alice, bob)
		}
	}
}

func condVariable() {
	count := 0
	finished := 0
	var mu sync.Mutex

	// this is the conditional variable
	cond := sync.NewCond(&mu)

	go func() {
		mu.Lock()
		defer mu.Unlock()
		// if a certain condition true
		if count%2 == 0 {
			count++
		}
		finished++
		cond.Broadcast()
	}()

	mu.Lock()
	defer mu.Unlock()

	// we're waiting for a shared variable to reach a certain condition
	for count < 5 && finished < 10 {
		cond.Wait()
	}

	if count >= 5 {
		fmt.Println("received 5+ votes!")
	} else {
		fmt.Println("lost.")
	}
}

func simpleChannel() {
	ch := make(chan bool)

	go func() {
		time.Sleep(1 * time.Second)
		<-ch
	}()

	start := time.Now()

	ch <- true

	fmt.Printf("send took %v\n", time.Since(start))
}

func rpcWithChannels() {
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func(x int) {
			rpcFunc(x)

			// write into the channel
			done <- true
		}(i)
	}

	for i := 0; i < 5; i++ {
		// read from the channel concurrently
		<-done
	}

}

func rpcFunc(num int) {
	fmt.Println(num)
}
