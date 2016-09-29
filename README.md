## Async channels 

This package implements some improved channels to golang. 

_Main problem_: Golang has no native concurent queues. But sometimes you have api or other microservice that reacts too long. 
"Сlassic" way to solve problem is channels and buffer channels. But buffer channels is limited. With highload services you often cant expect exact number of connections, requests, workers etc. 

I was inspired by [this article](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/) and desided to make something like proxy or async channels that can work with unlimited buffer using native language tools without extra locks. 


### Variables of proxying

#### 1. You workers are fast

In this case receiving item will be immidiately placed to job channel:
```
1/received message: 1-test11
1/received message (1-test11) proxified to job, 0 latency
1/received message: 1-test12
1/received message (1-test12) proxified to job, 0 latency
Working on "1-test12"...
Working on "1-test11"...
1/received message: 1-test13
1/received message (1-test13) proxified to job, 0 latency
1/received message: 1-test14
1/received message (1-test14) proxified to job, 0 latency
Working on "1-test14"...
Working on "1-test13"...
```

#### 1. You workers are slow

In this case, if nobody is listening to job channel, you items will be placed to temp slice, and send to job channel in future.
```
1/received message: 2-test9
1/received message (2-test9) proxified to job, 0 latency
Working on "2-test9"...
Working on "2-test8"...
1/received message: 2-test10
2/received message: 2-test11, push to long list
2/received message: 2-test12, push to long list
2/received message: 2-test13, push to long list
...
2/send (2-test13) async to job...
2/send (2-test12) async to job...
2/send (2-test11) async to job...
2/send (2-test10) async to job...
...
Working on "2-test13"...
Working on "2-test12"...
Working on "2-test11"...
```

