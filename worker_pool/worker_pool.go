package main

import (
	"fmt"
	"time"
)

const MaxGoroutines = 5

type WorkerPool struct {
	maxSize int // The limit size of the Goroutines
	pools chan struct{} // Pool queue
}

/*
 * @brief: Initialize the Goroutine pool
 * @params:
 *		maxSize: the limit size
 */
func (wp *WorkerPool)Init(maxSize int) {
	wp.maxSize = maxSize
	wp.pools = make(chan struct{}, maxSize)
}

/*
 * @brief: Add a function to the pool
 * @params:
 *		f func(input int) int: a function takes an input as an integer and returns an integer
 * 		input int: input value of the given function
 *		cb func(input int, output int): the callback function of an input value and returns the output value
 */
func (wp *WorkerPool)Add(f func(input int) int,
	input int, cb func(input int, output int)){
	wp.pools <- struct{}{}
	go wp.executeTask(f, input, cb) // Run the function in a new Goroutine
}

/*
 * @brief: Execute task in a Goroutine
 * @params:
 *		f func(input int) int: a function takes an input as an integer and returns an integer
 * 		input int: input value of the given function
 *		cb func(input int, output int): the callback function of an input value and returns the output value
 */
func (wp *WorkerPool)executeTask(f func(input int) int,
	input int, cb func(input int, output int)) {
	ret := f(input) // Execute function
	<- wp.pools
	cb(input, ret) // Callback function
}

func task(in int) int {
	time.Sleep(time.Second)
	fmt.Printf("===> Doing task %d\n", in)
	return in*2
}

func main() {
	var wp WorkerPool
	wp.Init(MaxGoroutines)

	for i := 0; i < 20; i++ {
		wp.Add(task, i, func(in int, out int) {
			fmt.Printf(">>>> Task id %d returns %d\n", in, out)
		})
	}

	time.Sleep(2*time.Second)
}