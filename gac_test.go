package gac

import (
	"fmt"
	"sync/atomic"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

const ITEMS_LEN = 1000000

func BenchmarkBuffChan(b *testing.B) {

	var bchan = make(chan int, ITEMS_LEN)

	for i := 0; i < ITEMS_LEN; i++ {
		bchan <- i
	}
}

func BenchmarkAsyncChanPooled(b *testing.B) {

	f, _ := NewAsyncChannel()

	for i := 0; i < ITEMS_LEN; i++ {
		f <- i
	}
}

func workerFast(pool chan interface{}, count *uint64) {
	for {
		<-pool
		// tmp := <-pool
		atomic.AddUint64(count, 1)
		// fmt.Printf("Working on %q...\n", tmp.(string))
	}
}

func workerSlow(pool chan interface{}, count *uint64, d time.Duration) {
	for {
		<-pool
		// tmp := <-pool
		atomic.AddUint64(count, 1)
		// fmt.Printf("Working on %q...\n", tmp.(string))
		time.Sleep(d)
	}
}

func TestFastPool(t *testing.T) {

	var ops uint64

	// Send data to receive channel, read from pool
	receive, pool := NewAsyncChannel()

	// spawn N fast workers
	for i := 0; i < 10; i++ {
		go workerFast(pool, &ops)
	}

	// will push data to channel with no locks and timeouts
	for i := 0; i < 100; i++ {
		receive <- fmt.Sprintf("1-test%d", i)
	}

	// tests fix for goroutines
	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, uint64(100), atomic.LoadUint64(&ops))

}

func TestSlowPool(t *testing.T) {

	var ops uint64

	// Send data to receive channel, read from pool
	receive, pool := NewAsyncChannel()

	// spawn N slow workers
	for i := 0; i < 10; i++ {
		go workerSlow(pool, &ops, 100*time.Millisecond)
	}

	// send first part of data to channel (will push data to channel with no locks and timeouts)
	for i := 0; i < 100; i++ {
		receive <- fmt.Sprintf("1-test%d", i)
	}

	// long pool stacking
	assert.NotEqual(t, uint64(100), atomic.LoadUint64(&ops))

	// wait 1 sec for all workers finish + 5 millisecond for test fix
	time.Sleep(1005 * time.Millisecond)

	// now all tasks must be finished
	assert.Equal(t, uint64(100), atomic.LoadUint64(&ops))

}
