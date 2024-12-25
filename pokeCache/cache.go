package pokeCache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	value     string
}

type Cache struct {
	cache    map[int]CacheEntry
	interval time.Duration
	lock     sync.RWMutex
}

func NewCache(reapInterval time.Duration) *Cache {
	cachMap := make(map[int]CacheEntry)
	cache := Cache{cache: cachMap, interval: reapInterval}
	go cache.reapLoop()
	return &cache
}

func (C *Cache) Add(key int, value string) {
	C.lock.Lock()
	cacheEntry := CacheEntry{time.Now(), value}
	C.cache[key] = cacheEntry
	C.lock.Unlock()
}

func (C *Cache) AddAll(firstId int, values []string) {
	for index, value := range values {
		C.Add(firstId+index, value)
	}
}

func (C *Cache) Get(key int) (string, bool) {
	C.lock.RLock()
	defer C.lock.RUnlock()
	value, ok := C.cache[key]
	return value.value, ok
}

func (C *Cache) GetRange(firstId, lastId int) []string {
	data := make([]string, lastId-firstId+1)
	C.lock.RLock()
	defer C.lock.RUnlock()
	for id := firstId; id < lastId+1; id++ {
		datum, ok := C.cache[id]
		if !ok {
			return nil
		}
		data[id-firstId] = datum.value
	}
	return data
}

func (C *Cache) reapLoop() {
	ticker := time.NewTicker(C.interval)
	for {
		<-ticker.C
		C.lock.Lock()
		for key, entry := range C.cache {
			if time.Since(entry.createdAt) > C.interval {
				delete(C.cache, key)
			}
		}
		C.lock.Unlock()
	}
}
