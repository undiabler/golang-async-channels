// -------------------------
// Copyright 2016, undiabler
// git: github.com/undiabler/golang-async-channels
// http://undiabler.com
//--------------------------

package main

import (
	"fmt"
	"time"
)

type proxy_tube struct {

	chan_from chan interface{}
	chan_to   chan interface{}

}

func NewAsyncChannel() (chan_from,chan_to chan interface{}) {
	
	//TODO: think about returning proxy_tube struct to avoid memory leaks
	p := new(proxy_tube)

	//TODO: maybe sometimes you will want buffered channels for even more amortization
	p.chan_from = make(chan interface{})
	p.chan_to = make(chan interface{})

	go p.proxy_worker()

	return p.chan_from,p.chan_to
}

func (p *proxy_tube) proxy_worker() {

	var item interface{}

	var items []interface{}

	for {

		if item == nil {

			select {

				case tmp := <- p.chan_from:

			        fmt.Printf("1/received message: %s\n", tmp)		

			        select {
				        case p.chan_to <- tmp:

					        fmt.Printf("1/received message (%s) proxified to job, 0 latency\n", tmp)

					        continue
					    default:
					    	item = tmp
			        }	
		    }

		} else {

			select {

			    case tmp := <- p.chan_from:

			        fmt.Printf("2/received message: %s, push to long list\n", tmp)

			        items = append(items,item)
			        item = tmp

			    case p.chan_to <- item:

			        fmt.Printf("2/send (%s) async to job...\n", item)

			    	item = nil

			        ln := len(items)

			        if ln > 0 {
			        	item = items[ln-1]
			        	items = items[:ln-1]
			        }
			        
		    }

		}


	}
}

const (
	WORKERS = 10
	WORKER_SLEEP = 50
)

func worker(pool chan interface{}) {
	for {
		tmp := <- pool
		fmt.Printf("Working on %q...\n", tmp.(string))
		time.Sleep( WORKER_SLEEP*time.Millisecond )
	}
}

func main() {

	receive, pool := NewAsyncChannel()

    for i := 0; i < WORKERS; i++ {
	    go worker(pool)
    }

   	for i := 0; i < 100; i++ {
   		receive <- fmt.Sprintf("1-test%d",i)
   	}

   	time.Sleep(1*time.Second)

   	for i := 0; i < 50; i++ {
   		receive <- fmt.Sprintf("2-test%d",i)
   	}

   	time.Sleep(5*time.Second)
}