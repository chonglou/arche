package amqp

import (
	"fmt"
	"time"

	"github.com/chonglou/arche/web/queue"
	"github.com/streadway/amqp"
)

// New create a amqp queue
func New(url, name string) queue.Queue {
	return &Queue{
		url:  url,
		name: name,
	}
}

// Queue queue for amqp
type Queue struct {
	url  string
	name string
}

// Put send a message
func (p *Queue) Put(typ, id string, pri uint8, buf []byte) error {
	return p.open(func(ch *amqp.Channel) error {
		return ch.Publish("", p.name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId:    id,
			Priority:     pri,
			Body:         buf,
			Timestamp:    time.Now(),
			Type:         typ,
		})
	})
}

// Launch launch a worker
func (p *Queue) Launch(name string) error {
	return p.open(func(ch *amqp.Channel) error {
		if err := ch.Qos(1, 0, false); err != nil {
			return err
		}
		msgs, err := ch.Consume(p.name, name, false, false, false, false, nil)
		if err != nil {
			return err
		}
		for d := range msgs {
			d.Ack(false)
			hnd := queue.Get(d.Type)
			if hnd == nil {
				return fmt.Errorf("unknown message type %s", d.Type)
			}
			if err := hnd(d.MessageId, d.Body); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Queue) open(f func(*amqp.Channel) error) error {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return f(ch)
}
