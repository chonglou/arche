package amqp

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/chonglou/arche/web/queue"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// New create a amqp queue
func New(url, name string) queue.Queue {
	return &Queue{
		url:       url,
		name:      name,
		consumers: make(map[string]queue.Consumer),
	}
}

// Queue queue for amqp
type Queue struct {
	url       string
	name      string
	consumers map[string]queue.Consumer
}

// Status status
func (p *Queue) Status() (map[string]interface{}, error) {
	tasks := make(map[string]string)
	for n, f := range p.consumers {
		tasks[n] = runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	}
	return map[string]interface{}{
		"tasks": tasks,
	}, nil
}

// Register register handler
func (p *Queue) Register(n string, c queue.Consumer) {
	if _, ok := p.consumers[n]; ok {
		log.Errorf("consumer for %s already exists, will override it", n)
	}
	p.consumers[n] = c
}

// Put send a message
func (p *Queue) Put(typ, id string, pri uint8, buf []byte) error {
	log.Debugf("send message %s@%s", id, typ)
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
		if _, err := ch.QueueDeclare(
			p.name, // name of the queue
			true,   // durable
			false,  // delete when unused
			false,  // exclusive
			false,  // noWait
			nil,    // arguments
		); err != nil {
			return err
		}
		msgs, err := ch.Consume(p.name, name, false, false, false, false, nil)
		if err != nil {
			return err
		}
		for d := range msgs {
			d.Ack(false)
			log.Debugf("receive message %s@%s", d.MessageId, d.Body)
			hnd, ok := p.consumers[d.Type]
			if !ok {
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
