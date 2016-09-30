package gac

import "testing"

const ITEMS_LEN = 10000

func BenchmarkBuffChan(b *testing.B) {

	var bchan = make(chan int,ITEMS_LEN)

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
