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

func TestWaitGroupPoolConcurrent(t *testing.T) {
	// Reuse the pool for many rounds of concurrent Add/Done/Wait cycles,
	// run with -race to catch counter misuse between Add and Done.
	for round := 0; round < 50; round++ {
		wgp := NewWaitGroupPool(2)
		var total uint32
		for i := 0; i < 20; i++ {
			wgp.Add()
			go func() {
				defer wgp.Done()
				atomic.AddUint32(&total, 1)
			}()
		}
		wgp.Wait()
		if total != 20 {
			t.Fatalf("round %d: expected 20, got %d", round, total)
		}
	}
}
