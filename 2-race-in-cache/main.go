package main

import (
	"container/list"
	"sync"
)

// CacheSize determines how big the cache can grow
const CacheSize = 100

// KeyStoreCacheLoader is an interface for the KeyStoreCache
type KeyStoreCacheLoader interface {
	// Load implements a function where the cache should get its content from
	Load(string) string
}

type page struct {
	Key   string
	Value string
}

// KeyStoreCache is a LRU cache for string key-value pairs
type KeyStoreCache struct {
	lock  sync.RWMutex // Use RWMutex for thread safety
	cache map[string]*list.Element
	pages list.List
	load  func(string) string
}

// New creates a new KeyStoreCache
func New(load KeyStoreCacheLoader) *KeyStoreCache {
	return &KeyStoreCache{
		load:  load.Load,
		cache: make(map[string]*list.Element),
	}
}

// Get gets the key from cache, loads it from the source if needed
func (k *KeyStoreCache) Get(key string) string {
	k.lock.RLock() // Acquire a read lock
	if e, ok := k.cache[key]; ok {
		// We've already read from the cache, release the read lock.
		// We'll now be able to acquire a write lock and move the element to the front of the list
		k.lock.RUnlock()
		k.lock.Lock()
		k.pages.MoveToFront(e)
		value := e.Value.(page).Value
		k.lock.Unlock()
		return value
	}
	k.lock.RUnlock() // Release the read lock

	k.lock.Lock() // Acquire a write lock

	// Check again in case the value was loaded by another writer while waiting for the lock
	if e, ok := k.cache[key]; ok {
		k.pages.MoveToFront(e)
		value := e.Value.(page).Value
		k.lock.Unlock() // Release the write lock
		return value
	}

	// Miss - load from database and save it in cache
	p := page{key, k.load(key)}
	// if cache is full, remove the least used item
	if len(k.cache) >= CacheSize {
		end := k.pages.Back()
		// remove from map
		delete(k.cache, end.Value.(page).Key)
		// remove from list
		k.pages.Remove(end)
	}
	k.pages.PushFront(p)
	k.cache[key] = k.pages.Front()
	value := p.Value

	k.lock.Unlock() // Release the write lock

	return value
}

// Loader implements KeyStoreCacheLoader
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

func run() *KeyStoreCache {
	loader := Loader{
		DB: GetMockDB(),
	}
	cache := New(&loader)

	RunMockServer(cache)

	return cache
}

func main() {
	run()
}
