package cache

import (
	"sync"
	"time"
)

type MemoryCachedBytes struct {
	value             []byte
	expireAtTimestamp time.Time
}

type MemoryCache struct {
	stopChan chan struct{}

	wg          sync.WaitGroup
	mu          sync.RWMutex
	cachedBytes map[string]MemoryCachedBytes
}

func (mc *MemoryCache) Get(key string) ([]byte, bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if v, ok := mc.cachedBytes[key]; ok {
		return v.value, true, nil
	}
	return []byte{}, false, nil
}

func (mc *MemoryCache) Put(key string, value []byte, lifetime time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cachedBytes[key] = MemoryCachedBytes{
		value:             value,
		expireAtTimestamp: time.Now().Add(lifetime),
	}
	return nil
}

func (mc *MemoryCache) Delete(key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.cachedBytes, key)
	return nil
}

func (mc *MemoryCache) Exists(key string) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	_, exists := mc.cachedBytes[key]
	return exists, nil
}

func (mc *MemoryCache) start(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-mc.stopChan:
			return
		case <-t.C:
			mc.mu.RLock()
			for k := range mc.cachedBytes {
				if time.Now().After(mc.cachedBytes[k].expireAtTimestamp) {
					mc.mu.RUnlock()
					mc.Delete(k)
					mc.mu.RLock()
				}
			}
			mc.mu.RUnlock()
		}
	}
}

func (mc *MemoryCache) stop() {
	close(mc.stopChan)
	mc.wg.Wait()
}

func NewMemoryCache(cleanupInterval time.Duration) *MemoryCache {
	mc := &MemoryCache{
		cachedBytes: make(map[string]MemoryCachedBytes),
		stopChan:    make(chan struct{}),
	}

	mc.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer mc.wg.Done()
		mc.start(cleanupInterval)
	}(cleanupInterval)

	return mc
}
