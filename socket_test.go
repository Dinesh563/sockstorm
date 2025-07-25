package main

import (
	"sync"
	"testing"
)

func TestWith10Connections(t *testing.T) {
	var wg sync.WaitGroup
	

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Test1()
		}()
	}

	wg.Wait()
}
