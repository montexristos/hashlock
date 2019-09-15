package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	l := (&LockHash{}).New()
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go iteration(l, i, t, &wg)
	}
	wg.Wait()
}

func iteration(l LockHash, i int, t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(10-2*i) * time.Second)
	t.Log(fmt.Sprintf("get Lock %d", i))
	l.Lock("key")
	l.Unlock("key")
	t.Log(fmt.Sprintf("unLock %d", i))
}
