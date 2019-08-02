package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/theTardigrade/fan"
)

type fibonacciCacheDatum struct {
	data  map[int]uint64
	mutex sync.RWMutex
}

func newFibonacciCacheDatum() *fibonacciCacheDatum {
	return &fibonacciCacheDatum{
		data: make(map[int]uint64),
	}
}

func (d *fibonacciCacheDatum) Get(n int) (value uint64, found bool) {
	d.mutex.RLock()
	value, found = d.data[n]
	d.mutex.RUnlock()
	return
}

func (d *fibonacciCacheDatum) Set(n int, value uint64) {
	d.mutex.Lock()
	d.data[n] = value
	d.mutex.Unlock()
}

var (
	fibonacciCache = newFibonacciCacheDatum()
)

func fibonacci(n int) uint64 {
	if n == 0 || n == 1 {
		return uint64(n)
	}

	if value, found := fibonacciCache.Get(n); found {
		return value
	}

	value := fibonacci(n-1) + fibonacci(n-2)
	fibonacciCache.Set(n, value)
	return value
}

func fibonacciHandler(n int) error {
	if n > 0 && n%66 == 0 {
		return errors.New("cannot be divisible by sixty-six")
	}

	fmt.Printf("%d = %d\n", n, fibonacci(n))

	return nil
}

func main() {
	err := fan.HandleRepeated(fibonacciHandler, 1e12)
	fmt.Println("ERR", err)
}
