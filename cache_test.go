package cache_test

import (
	"testing"
	"time"

	cache "github.com/evamayerova/lrucache"
)

type CacheItems []interface{}

func TestCache(t *testing.T) {
	c := cache.NewCache(5, "0")
	items := CacheItems{}
	for i := 0; i < 10; i++ {
		items = append(items, i)
	}

	for _, k := range items {
		c.Put(k, 1, 300, 1)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5 * time.Millisecond)
	for _, k := range items[:5] {
		if c.Read(k) != nil {
			t.Errorf("Item found in cache")
		}
	}

	for _, k := range items[5:] {
		if c.Read(k) == nil {
			t.Errorf("Item not found in cache")
		}
	}
}

func TestChance(t *testing.T) {
	c := cache.NewCache(5, "0")
	c.Put("0", 1, 300, 0)
	time.Sleep(5 * time.Millisecond)
	if c.Read("0") != nil {
		t.Errorf("Element had 0 chance to be put, but still was")
	}
	c.Put("0", 1, 300, 1)
	time.Sleep(5 * time.Millisecond)
	if c.Read("0") == nil {
		t.Errorf("Element had 1 chance to be put, but wasnt")
	}
}

func TestTTL(t *testing.T) {
	c := cache.NewCache(5, "0")
	var keys []int
	for i := 0; i < 5; i++ {
		keys = append(keys, i)
	}

	c.Put(keys[0], 0, 1, 1)
	time.Sleep(2000 * time.Millisecond)
	if c.Read(keys[0]) != nil {
		t.Error("TTL does not work")
	}

	for i := 1; i < 5; i++ {
		c.Put(keys[i], 0, 1, 1)
	}

	time.Sleep(2000 * time.Millisecond)
	for i := 0; i < 5; i++ {
		if c.Read(keys[i]) != nil {
			t.Error("TTL does not work")
		}
	}
}
