package lrucache

import (
	"fmt"
	"sync"
	"time"
)

// Item specified by key, value and ttl
type Item struct {
	key      interface{}
	value    interface{}
	deadline int64
}

// NewItem creates new instance of Item
func NewItem(k, value interface{}, ttl int) *Item {
	now := time.Now().Unix()
	return &Item{
		key:      k,
		value:    value,
		deadline: now + int64(ttl),
	}
}

// Cache structure
type Cache struct {
	newBuffer map[interface{}]Item
	oldBuffer map[interface{}]Item
	capacity  int
	mtx       sync.RWMutex
	inBuffer  chan *Item
}

// NewCache creates new instance of Cache
func NewCache(cap int) *Cache {
	c := &Cache{
		newBuffer: map[interface{}]Item{},
		oldBuffer: map[interface{}]Item{},
		capacity:  cap,
		inBuffer:  make(chan *Item, cap),
	}

	// cache Writer
	go func() {
		var i *Item
		for {
			i = <-c.inBuffer
			c.put(i)
		}
	}()
	return c
}

// Len returns length of cache buffer
func (c *Cache) Len() int {
	return len(c.newBuffer)
}

// Read cached item if present, otherwise return nil
func (c *Cache) Read(key interface{}) interface{} {
	c.rLockMtx()
	if item, ok := c.newBuffer[key]; ok {
		c.rUnlockMtx()
		// ttl
		if item.deadline < time.Now().Unix() {
			return nil
		}
		return item.value
	}
	if item, ok := c.oldBuffer[key]; ok {
		c.rUnlockMtx()
		// ttl
		if item.deadline < time.Now().Unix() {
			return nil
		}
		c.asyncWrite(&item)
		return item.value
	}
	c.rUnlockMtx()
	return nil
}

// Write new item into cache. TTL specifies the maximum living time of a record.
func (c *Cache) Write(key, value interface{}, ttl int) error {
	if value == nil {
		return fmt.Errorf("tried to put nil value into cache")
	}
	ci := NewItem(key, value, ttl)
	c.asyncWrite(ci)
	return nil
}

func (c *Cache) put(i *Item) {
	// if item is already cached, update it's value
	if c.update(i.key, i.value) {
		return
	}
	//otherwise create new cache item
	c.write(i.key, i)
}

func (c *Cache) check(k interface{}) bool {
	c.rLockMtx()
	defer c.rUnlockMtx()
	if _, ok := c.newBuffer[k]; ok {
		return true
	}
	return false
}

func (c *Cache) update(k, v interface{}) bool {
	if !c.check(k) {
		return false
	}
	c.lockMtx()
	defer c.unlockMtx()
	if item, ok := c.newBuffer[k]; ok {
		item.value = v
		return true
	}
	return false
}

func (c *Cache) write(key interface{}, val *Item) {
	c.lockMtx()
	defer c.unlockMtx()
	c.newBuffer[key] = *val
	if c.Len() >= c.capacity {
		c.oldBuffer = c.newBuffer
		c.newBuffer = map[interface{}]Item{}
	}
}

func (c *Cache) asyncWrite(ci *Item) {
	select {
	case c.inBuffer <- ci:
	default:
		// cache write dropped
	}
}

func (c *Cache) lockMtx() {
	c.mtx.Lock()
}

func (c *Cache) unlockMtx() {
	c.mtx.Unlock()
}

func (c *Cache) rLockMtx() {
	c.mtx.RLock()
}

func (c *Cache) rUnlockMtx() {
	c.mtx.RUnlock()
}
