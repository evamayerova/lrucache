package cache_test

import (
	"testing"
	"time"

	cache "github.com/evamayerova/lrucache"
)

func TestCacheManager(t *testing.T) {
	c, _ := cache.NewManager(1, 5)

	for k := 0; k < 5; k++ {
		c.Put(k, 1, 300, 1)
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(5 * time.Millisecond)

	for k := 0; k < 5; k++ {
		if c.Read(k) == nil {
			t.Errorf("Item not found in cache")
		}
		time.Sleep(100 * time.Millisecond)
	}
	c.Read(0)
	c.Put(5, 1, 1, 1)
	time.Sleep(100 * time.Millisecond)
	c.Put(5, 1, 10, 1)
	time.Sleep(100 * time.Millisecond)
	c.Put(1, nil, 300, 1)
	time.Sleep(100 * time.Millisecond)
	err := c.Put(nil, 1, 300, 1)
	if err == nil {
		t.Errorf("nil key was put into cache successfully")
	}
	time.Sleep(time.Second * 2)
	if c.Read(5) != nil {
		t.Errorf("TTL not working")
	}
	if c.Read(1) == nil {
		t.Errorf("nil object was put into cache successfully")
	}
}

func TestDistribution(t *testing.T) {
	cm, _ := cache.NewManager(2, 3)

	cm.Put(0, 0, 300, 1)
	cm.Put(2, 0, 300, 1)
}
