# lrucache

Fast relaxed LRU implementation in Golang

## Implementation
The theorethical and explicit implementation of last-recently-used cache is expensive. It requires to have a structure (map) with the chached data and also a linked list which is ordered by the recency of the data writes or reads. The modification of the linked list is needed for every operation, and the locking of it slows the implementation.

This implementation relaxes the explicit order or the elements. That allows us to use RW locks, which lock the data structure only for writes. Instead of a map and a linked list, we use two maps, `old` and `new`. 

* The `old` map is used only for reads, thus it does not lock the structure for other goroutines. 
* The `new` map is used only for writes.
* When `new` map reaches the cache capacity, the items from `new` are moved to `old` with removing the previous `old` elements. `new` is initialized to empty map.
* First, we try to read from `new`. If we fail, and the item is in `old`, we move it from `old` to `new`. The moving happens through a channel, and then it is written to the map asynchronously in another routine, so the speed of the read is not affected.

With this approach, we can preserve the approximate LRU feature with keeping the read queries as fast as possible.


## Documentation
See [Godoc](https://godoc.org/github.com/evamayerova/lrucache)

## Cache
### Example usage
```
capacity := 10
c := lrucache.NewCache(capacity)

key := 0
val := 1
ttl := 300 // seconds
chance := float32(1) // 1.0 for 100% to be cached

c.Write(key, val, ttl, chance)
v := c.Read(key)
```

## Multiple caches

The lib also allows to create multiple caches and reduce the time on a mutex wait even more. With use of a cache `Manager`, you can specify the number of caches and it will distribute the keys evenly using modulo of the key. The condition is, that the keys must be convertable to `int64`. The interface is the same as for a single cache.

### Example usage
```
cacheNr := 2
capacity := 10 // overall capacity - it is splitted into the caches evenly
cm := lrucache.NewManager(cacheNr, capacity)

key := 0
val := 1
ttl := 300 // seconds
chance := float32(1) // 1.0 for 100% to be cached

cm.Write(key, val, ttl, chance)
v := cm.Read(key)
```
