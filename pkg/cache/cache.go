// Package cache menyediakan fungsionalitas cache dalam memori yang aman dari kondisi balapan (thread-safe).
package cache

import (
	"sync"
	"time"
)

var GlobalCache = KonstruksiCacheBaru(5*time.Minute, 1*time.Minute)

type Item struct {
	Value      interface{}
	Expiration int64
}

// Expired mengecek apakah item cache sudah kedaluwarsa.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type Cache struct {
	mu            sync.RWMutex
	items         map[string]Item
	defaultTTL    time.Duration
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

// KonstruksiCacheBaru membuat instance Cache in-memory thread-safe baru dengan TTL default dan interval pembersihan.
func KonstruksiCacheBaru(defaultTTL, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		items:       make(map[string]Item),
		defaultTTL:  defaultTTL,
		stopCleanup: make(chan struct{}),
	}
	
	if cleanupInterval > 0 {
		c.cleanupTicker = time.NewTicker(cleanupInterval)
		go c.startCleanupLoop()
	}
	
	return c
}

// Set menaruh item ke dalam cache dengan TTL tertentu.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	var expiration int64
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}
	
	c.mu.Lock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
	c.mu.Unlock()
}

// Get mengambil item dari cache, mengembalikan nil dan false jika tidak ditemukan atau expired.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()
	
	if !found {
		return nil, false
	}
	
	if item.Expired() {
		return nil, false
	}
	
	return item.Value, true
}

// Delete menghapus item tertentu dari cache berdasarkan key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Clear membersihkan seluruh isi cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	c.items = make(map[string]Item)
	c.mu.Unlock()
}

// Close menghentikan background pembersihan cache.
func (c *Cache) Close() {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
		close(c.stopCleanup)
	}
}

func (c *Cache) startCleanupLoop() {
	for {
		select {
		case <-c.cleanupTicker.C:
			c.mu.Lock()
			now := time.Now().UnixNano()
			for k, v := range c.items {
				if v.Expiration > 0 && now > v.Expiration {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		case <-c.stopCleanup:
			return
		}
	}
}
