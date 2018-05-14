// -------------------------
// Copyright 2016, undiabler
// git: github.com/undiabler/golang-async-channels
// http://undiabler.com
//--------------------------

package gac

type proxyTube struct {
	chan_from chan interface{}
	chan_to   chan interface{}
}

func NewAsyncChannel() (chan_from, chan_to chan interface{}) {

	//TODO: think about returning proxyTube struct to avoid memory leaks
	p := new(proxyTube)

	//TODO: maybe sometimes you will want buffered channels for even more amortization
	p.chan_from = make(chan interface{})
	p.chan_to = make(chan interface{})

	go p.proxy_worker()

	return p.chan_from, p.chan_to
}

func (p *proxyTube) proxy_worker() {

	var items []interface{}

	for {

		items_len := len(items)

		if items_len == 0 {

			select {

			case tmp := <-p.chan_from:

				// fmt.Printf("1/received message: %s\n", tmp)

				select {

				case p.chan_to <- tmp:

					// fmt.Printf("1/received message (%s) proxified to job, 0 latency\n", tmp)

					continue

				default:
					items = append(items, tmp)
				}
			}

		} else {

			select {

			case tmp := <-p.chan_from:

				// fmt.Printf("2/received message: %s, push to long list\n", tmp)

				items = append(items, tmp)

			case p.chan_to <- items[items_len-1]:

				// fmt.Printf("2/send (%s) async to job...\n", items[items_len-1])

				items = items[:items_len-1]

			}

		}

	}
}
