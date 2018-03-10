package redis

import (
	"encoding/json"
	"fmt"

	"github.com/chonglou/arche/web/queue"
	_redis "github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// New create a amqp queue
func New(name string, pool *_redis.Pool) queue.Queue {
	return &Queue{
		name: name,
		pool: pool,
	}
}

// Queue by redis
type Queue struct {
	name string
	pool *_redis.Pool
}

// Put send a message
func (p *Queue) Put(t *queue.Task) error {
	log.Debugf("send message %s@%s", t.ID, t.Type)
	buf, err := json.Marshal(t)
	if err != nil {
		return err
	}
	_, err = p.pool.Get().Do("LPUSH", p.name, string(buf))
	return err
}

// Launch launch a worker
func (p *Queue) Launch(name string) error {
	res, err := _redis.Strings(p.pool.Get().Do("BRPOP", p.name))
	if err != nil {
		return err
	}
	var t queue.Task
	if err = json.Unmarshal([]byte(res[1]), &t); err != nil {
		return err
	}
	log.Debugf("receive message %s@%s", t.ID, t.Type)
	hnd, ok := queue.Get(t.Type)
	if !ok {
		return fmt.Errorf("unknown message type %s", t.Type)
	}
	return hnd(t.ID, t.Body)
}
