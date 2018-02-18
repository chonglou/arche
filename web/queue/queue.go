package queue

// Queue message queue
type Queue interface {
	Put(typ, id string, pri uint8, buf []byte) error
	Launch(name string) error
}
