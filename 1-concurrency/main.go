package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func main() {
	returnIntChan := make(chan int, 10)
	resultSlice := []int{}
	intSquare(returnIntChan)
	wg := sync.WaitGroup{}
	wg.Add(10)
	go func() {
		for num := range returnIntChan {
			resultSlice = append(resultSlice, num)
			wg.Done()
		}
	}()
	wg.Wait()
	fmt.Println(resultSlice)
}

func randomSlice(intChan chan int) {
	randSlice := []int{}
	wg := sync.WaitGroup{}
	wg.Add(10)
	go func() {
		for i := 0; i < 10; i++ {
			randSlice = append(randSlice, rand.IntN(101))
			wg.Done()
		}
	}()
	wg.Wait()
	for _, num := range randSlice {
		intChan <- num
	}
	close(intChan)
}

func intSquare(returnIntChan chan int) {
	intChan := make(chan int, 10)
	randomSlice(intChan)
	wg := sync.WaitGroup{}
	wg.Add(10)
	go func() {
		for num := range intChan {
			returnIntChan <- num * num
			wg.Done()
		}
	}()
	wg.Wait()
	close(returnIntChan)
}
