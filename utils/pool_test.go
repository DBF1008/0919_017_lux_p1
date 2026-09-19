package utils

import (
	"sync/atomic"
	"testing"
)

func TestWaitGroupPool(t *testing.T) {
	wgp := NewWaitGroupPool(10)

	var total uint32

	for i := 0; i < 100; i++ {
		wgp.Add()
		go func(total *uint32) {
			defer wgp.Done()
			atomic.AddUint32(total, 1)
		}(&total)
	}
	wgp.Wait()

	if total != 100 {
		t.Fatalf("The size '%d' of the pool did not meet expectations.", total)
	}
}

func TestWaitGroupPoolConcurrentAddDone(t *testing.T) {
	// Stress Add/Done/Wait under high concurrency to verify there is no
	// data race or WaitGroup misuse panic (run with -race).
	for round := 0; round < 50; round++ {
		wgp := NewWaitGroupPool(4)
		var total int32
		for i := 0; i < 200; i++ {
			wgp.Add()
			go func() {
				defer wgp.Done()
				atomic.AddInt32(&total, 1)
			}()
		}
		wgp.Wait()
		if total != 200 {
			t.Fatalf("round %d: expected 200, got %d", round, total)
		}
	}
}
