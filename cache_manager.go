package cache

import (
	"fmt"

	"github.com/evamayerova/lrucache/utils"
)

// Manager serves as a cache sharder - it initiates multiple caches (based on `cacheCnt` parameter) and distributes incomming records evenly between them. All caches share the same maximal capacity defined by param `cap`
type Manager struct {
	caches   map[int]*Cache
	capacity int
	cacheCnt int
}

// NewManager creates a new CacheManager instance.
func NewManager(cacheCnt int, cap int) (*Manager, error) {
	if cacheCnt <= 0 {
		return nil, fmt.Errorf("cache count must be larger than 0")
	}
	cacheCap := cap / cacheCnt
	cm := Manager{
		capacity: cap,
		cacheCnt: cacheCnt,
	}
	cm.caches = make(map[int]*Cache, cacheCnt)
	for i := 0; i < cacheCnt; i++ {
		cm.caches[i] = NewCache(cacheCap)
	}
	return &cm, nil
}

func (cm *Manager) selectCache(k int64) (*Cache, error) {
	if k < 0 {
		k = -k
	}
	return cm.caches[int(k)%cm.cacheCnt], nil
}

// Put new item into cache. The key must be a numerical interface (convertible to int64), otherwise it will lead to an error.
func (cm *Manager) Put(k, v interface{}, ttl int, chance float32) error {
	num, err := utils.NumericInterfToInt(k)
	if err != nil {
		return err
	}
	cache, err := cm.selectCache(num)
	if err == nil {
		return cache.Put(k, v, ttl, chance)
	}
	return err
}

// Read item from cache. If the key was not found in cache or the key is not convertible to int64, it will return nil.
func (cm *Manager) Read(k interface{}) interface{} {
	num, err := utils.NumericInterfToInt(k)
	if err != nil {
		return nil
	}
	cache, err := cm.selectCache(num)
	if err == nil {
		return cache.Read(k)
	}
	return nil
}
