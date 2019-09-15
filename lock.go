package main

import (
	"strings"
	"sync"
	"time"
)

type LockHash struct {
	locks map[string]*sync.Mutex
}

func (l LockHash) New() LockHash {
	l.locks = make(map[string]*sync.Mutex)
	return l
}

func (l LockHash) GetLock(key string) *sync.Mutex {
	var found bool
	if _, found = l.locks[key]; !found {
		l.locks[key] = &sync.Mutex{}
	}
	return l.locks[key]
}

func (l LockHash) Lock(key string) {
	l.GetLock(key).Lock()
	time.Sleep(1 * time.Second)
}

func (l LockHash) Unlock(key string) {
	l.GetLock(key).Unlock()
}

func (l LockHash) GetLockKey(args []string) string {
	return strings.Join(args, "-")
}
