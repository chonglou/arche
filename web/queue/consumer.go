package queue

import "fmt"

// Consumer task handler
type Consumer func(id string, body []byte) error

var consumers = make(map[string]Consumer)

// Register register handler
func Register(n string, c Consumer) {
	if _, ok := consumers[n]; ok {
		panic(fmt.Sprintf("consumer for %s already exists, will override it", n))
	}
	consumers[n] = c
}

// Get get consumer by name
func Get(n string) Consumer {
	return consumers[n]
}
