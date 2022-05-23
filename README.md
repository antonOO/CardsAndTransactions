### Cards and transactions solution 

The solution contains a module transaction, which introduces a Card and 
Client data structures. 

All Client methods are synced and thread-safe. The method for retrieving 
cards is also guarded, but reads may be concurrent. 


#### Installation and testing guide

- Prerequisites - installed golang 1.17

1. Clone the repo 
2. cd into CardsAndTransactions - `go build ./...`
3. cd into /transactions - `go test ./...`

Demo live: 
repeat 1. and 2. 
3. run the test server - `go run main.go`
4. Experiment with requests similar to: 

```
$ curl -X POST localhost:8081/add-card -d '{"id": 1234567890123456}'
Added

$ curl -X POST localhost:8081/add-card -d '{"id": 1234567890123456}'
card with Id - 1234567890123456, already exists

$ curl -X POST localhost:8081/activate -d '{"id": 1234567890123456}'
Activated

$ curl -X GET localhost:8081/cards 
ID - 1234567890123456, Last Used - 2022-05-23 06:22:56.667378903 -0700 PDT m=+20.082500656 
^ ACTIVE CARDS
```

NOTES: 
Alternative solution is using a priority queue for disabling last activated cards
However, with the cost of card addition time.