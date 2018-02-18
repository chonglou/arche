package cache_test

import (
	"testing"
	"time"

	"github.com/chonglou/arche/web/cache"
	"github.com/chonglou/arche/web/cache/redis"
)

type M struct {
	S string
	I int
	T time.Time
}

func TestRedis(t *testing.T) {
	testCache(t, redis.New("redis://localhost:6379/6", "cache://"))
}

func testCache(t *testing.T, ch cache.Cache) {
	m1 := M{
		S: "hello, cache",
		I: 123,
		T: time.Now(),
	}
	t.Logf("%+v", m1)
	const key = "hi"
	if err := ch.Put(key, &m1, time.Hour); err != nil {
		t.Fatal(err)
	}

	var m2 M
	if err := ch.Get(key, &m2); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", m2)
	if m1.S != m2.S || m1.I != m2.I {
		t.Fatalf("want %+v, get %+v", m1, m2)
	}

	ch.Put("test", 123, time.Hour)
	if err := ch.Clear(); err != nil {
		t.Fatal(err)
	}
}
