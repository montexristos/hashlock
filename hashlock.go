package main

import (
	"strings"
	"sync"
)

// HashLock keeps a hashmap with unique keys and lock conditions
// Also has a rw mutex for the hashmap to avoid race conditions between threads
type HashLock struct {
	locks   map[string]*sync.Cond
	mapLock *sync.RWMutex
}

// New initializes a new HashLock and returns the instance
func (l HashLock) New() HashLock {
	l.locks = make(map[string]*sync.Cond)
	l.mapLock = &sync.RWMutex{}
	return l
}

// GetLock returns a lock condition for the provided key
func (l HashLock) GetLock(key string) *sync.Cond {
	l.mapLock.RLock()
	_, found := l.locks[key]
	l.mapLock.RUnlock()
	if !found {
		l.mapLock.Lock()
		l.locks[key] = sync.NewCond(&sync.Mutex{})
		l.mapLock.Unlock()
	}
	return l.locks[key]
}

// Lock locks the provided key for rw
func (l HashLock) Lock(key string) {
	l.GetLock(key).L.Lock()
}

// Unlock unlocks the provided key for rw
func (l HashLock) Unlock(key string) {
	l.GetLock(key).Broadcast()
	l.GetLock(key).L.Unlock()
}

// GetLockKey provides a pattern for creating keys from string arrays
func (l HashLock) GetLockKey(args []string) string {
	return strings.Join(args, "-")
}

// DeleteKey removes a key from hashmap with locks
func (l HashLock) DeleteKey(key string) {
	var found bool
	if _, found = l.locks[key]; !found {
		l.mapLock.Lock()
		defer l.mapLock.Unlock()
		delete(l.locks, key)
	}
}
