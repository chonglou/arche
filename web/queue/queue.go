package queue

// Queue message queue
type Queue interface {
	Put(*Task) error
	Launch(name string) error
}
