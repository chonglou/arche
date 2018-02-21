package queue

// Consumer task handler
type Consumer func(id string, body []byte) error

// Queue message queue
type Queue interface {
	Register(n string, c Consumer)
	Put(typ, id string, pri uint8, buf []byte) error
	Launch(name string) error
}
