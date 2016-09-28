// -------------------------
// Copyright 2015, undiabler
// git: github.com/undiabler/golang-async-channels
// http://undiabler.com
//--------------------------

package main

import (
	"fmt"
	"time"
)

var receive = make(chan string)
var pool = make(chan string)


func pooler(messages chan string, job chan string) {

	var item string

	var items []string

	for {

		if item == "" {

			select {

				case tmp := <- messages:

			        fmt.Printf("1/received message: %s\n", tmp)		

			        select {
				        case job <- tmp:
					        fmt.Printf("1/received message (%s) proxified to job, 0 latency\n", tmp)		
					    default:
					    	item = tmp
			        }	
		    }

		} else {

			select {

			    case tmp := <- messages:

			        fmt.Printf("2/received message: %s, push to long list\n", tmp)

			        items = append(items,item)
			        item = tmp

			    case job <- item:

			        fmt.Printf("2/send (%s) async to job...\n", item)

			    	item = ""

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

func worker() {
	for {
		tmp := <- pool
		fmt.Printf("Working on %q...\n", tmp)
		time.Sleep( WORKER_SLEEP*time.Millisecond )
	}
}

func main() {

    go pooler(receive, pool)

    for i := 0; i < WORKERS; i++ {
	    go worker()
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