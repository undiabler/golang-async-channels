## Async channels 

This package implements some improved channels to golang. 

_Main problem_: Golang has no native concurent queues. But sometimes you have api or other microservice that reacts too long. 
"Сlassic" way to solve problem is channels and buffer channels. But buffer channels is limited. With highload services you often cant expect exact number of connections, requests, workers etc. 

I was inspired by [this article](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/) and desided to make something like proxy or async channels that can work with unlimited buffer using native language tools without extra locks. 