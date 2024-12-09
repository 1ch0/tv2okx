package utils

import (
	"context"
	"sync"
	"time"
)

// memoryCache memory cache, support time expired
type memoryCache struct {
	data          interface{}
	cacheDuration time.Duration
	startTime     time.Time
}

// NewMemoryCache new memory cache instance
func newMemoryCache(data interface{}, cacheDuration time.Duration) *memoryCache {
	mc := &memoryCache{data: data, cacheDuration: cacheDuration, startTime: time.Now()}
	return mc
}

// IsExpired whether the cache data expires
func (m *memoryCache) IsExpired() bool {
	if m.cacheDuration <= 0 {
		return false
	}
	return time.Now().After(m.startTime.Add(m.cacheDuration))
}

// GetData get cache data
func (m *memoryCache) GetData() interface{} {
	return m.data
}

// MemoryCacheStore a sample memory cache instance, if data set cache duration, will auto clear after timeout.
// But, Expired cleanup is not necessarily accurate, it has a 3-second window.
type MemoryCacheStore struct {
	store sync.Map
}

// NewMemoryCacheStore memory cache store
func NewMemoryCacheStore(ctx context.Context) *MemoryCacheStore {
	mcs := &MemoryCacheStore{
		store: sync.Map{},
	}
	go mcs.run(ctx)
	return mcs
}

func (m *MemoryCacheStore) run(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.store.Range(func(key, value interface{}) bool {
				if value.(*memoryCache).IsExpired() {
					m.store.Delete(key)
				}
				return true
			})
		}
	}
}

// Put cache data, if cacheDuration>0, store will clear data after timeout.
func (m *MemoryCacheStore) Put(key, value interface{}, cacheDuration time.Duration) {
	mc := newMemoryCache(value, cacheDuration)
	m.store.Store(key, mc)
}

// Delete cache data from store
func (m *MemoryCacheStore) Delete(key interface{}) {
	m.store.Delete(key)
}

// Get cache data from store, if not exist or timeout, will return nil
func (m *MemoryCacheStore) Get(key interface{}) (value interface{}) {
	mc, ok := m.store.Load(key)
	if ok && !mc.(*memoryCache).IsExpired() {
		return mc.(*memoryCache).GetData()
	}
	return nil
}
