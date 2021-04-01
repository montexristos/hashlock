package hashlock

import (
	"strings"
	"sync"
	"time"
)

// HashLock keeps a hashmap with unique keys and lock conditions
// Also has a rw mutex for the hashmap to avoid race conditions between threads
// If timeout is set the lock is unlocked after the given time
type HashLock struct {
	locks   map[string]chan bool
	mapLock *sync.RWMutex
	timeout time.Duration
}

// New initializes a new HashLock
func (l *HashLock) New(d time.Duration) *HashLock {
	l.locks = make(map[string]chan bool)
	l.mapLock = &sync.RWMutex{}
	if d > 0 {
		l.timeout = d
	} else {
		l.timeout = 1 * time.Second
	}
	return l
}

// GetLock returns a lock condition for the provided key
func (l *HashLock) GetLock(key string) chan bool {
	l.mapLock.RLock()
	_, found := l.locks[key]
	if !found {
		//remove the read lock and lock for writes
		l.mapLock.RUnlock()
		l.mapLock.Lock()
		defer l.mapLock.Unlock()
		l.locks[key] = make(chan bool, 1)
	} else {
		defer l.mapLock.RUnlock()
	}
	return l.locks[key]
}

// Lock locks the provided key for rw
func (l *HashLock) Lock(key string) {
	l.GetLock(key) <- true
	//unlock after one second no need to wait more
	if l.timeout > 0 {
		time.AfterFunc(1*l.timeout, func() {
			// loop to handle more than one locks
			for i := 0; i < len(l.GetLock(key)); i++ {
				<-l.GetLock(key)
			}
		})
	}
}

// Unlock unlocks the provided key for rw
func (l *HashLock) Unlock(key string) {
	if len(l.GetLock(key)) > 0 {
		<-l.GetLock(key)
	}
}

// GetLockKey provides a pattern for creating keys from string arrays
func (l *HashLock) GetLockKey(args []string) string {
	return strings.Join(args, "-")
}

// DeleteKey removes a key from hashmap with locks
func (l *HashLock) DeleteKey(key string) {
	var found bool
	if _, found = l.locks[key]; found {
		l.mapLock.Lock()
		defer l.mapLock.Unlock()
		delete(l.locks, key)
	}
}

// Empty removes all non used keys from hashmap to release memory
func (l *HashLock) Empty() bool {
	l.mapLock.Lock()
	defer l.mapLock.Unlock()
	for k, v := range l.locks {
		if len(v) == 0 {
			delete(l.locks, k)
		} else {
			return false
		}
	}
	return true
}
