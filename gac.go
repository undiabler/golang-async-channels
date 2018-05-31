// -------------------------
// Copyright 2016, undiabler
// git: github.com/undiabler/golang-async-channels
// http://undiabler.com
//--------------------------

package gac

type proxyTube struct {
	chanFrom chan interface{}
	chanTo   chan interface{}
}

// NewAsyncChannel creates input and output unlimited chans
func NewAsyncChannel() (chanFrom, chanTo chan interface{}) {

	//TODO: think about returning proxyTube struct to avoid memory leaks
	p := new(proxyTube)

	//TODO: maybe sometimes you will want buffered channels for even more amortization
	p.chanFrom = make(chan interface{})
	p.chanTo = make(chan interface{})

	go p.proxyWorker()

	return p.chanFrom, p.chanTo
}

func (p *proxyTube) proxyWorker() {

	var items []interface{}

	for {

		itemsLen := len(items)

		if itemsLen == 0 {

			select {

			case tmp := <-p.chanFrom:

				// fmt.Printf("1/received message: %s\n", tmp)

				select {

				case p.chanTo <- tmp:

					// fmt.Printf("1/received message (%s) proxified to job, 0 latency\n", tmp)

					continue

				default:
					items = append(items, tmp)
				}
			}

		} else {

			select {

			case tmp := <-p.chanFrom:

				// fmt.Printf("2/received message: %s, push to long list\n", tmp)

				items = append(items, tmp)

			case p.chanTo <- items[itemsLen-1]:

				// fmt.Printf("2/send (%s) async to job...\n", items[itemsLen-1])

				items = items[:itemsLen-1]

			}

		}

	}
}
