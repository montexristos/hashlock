package main

import (
	"sync"
	"testing"
	"time"
)

// TestDifferentKeys should check that locking and specific key in hash map does not block other lock operations
func TestDifferentKeys(t *testing.T) {
	hashLock := (HashLock{}).New()
	c1 := make(chan string, 2)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		key := hashLock.GetLockKey([]string{"test1"})
		hashLock.Lock(key)
		t.Log("Got lock 1")
		defer func() {
			t.Log("Released lock 1")
			c1 <- "result 1"
			wg.Done()
		}()
		time.Sleep(5 * time.Second)
	}()


	//make sure the 1st goroutine runs first
	time.Sleep(1 * time.Second)
	wg.Add(1)
	go func() {
		key := hashLock.GetLockKey([]string{"test2"})
		hashLock.Lock(key)
		t.Log("Got lock 2")
		defer func() {
			t.Log("Released lock 2")
			c1 <- "result 2"
			wg.Done()
		}()
		time.Sleep(1 * time.Second)
	}()
	wg.Wait()
	//get the first channel value, it should be result 1
	if value := <-c1 ; value != "result 2" {
		t.Error("1 returned first")
	}
	//get the latest channel value it should be result 2
	if value := <-c1 ; value != "result 1" {
		t.Error("2 returned last")
	}
	close(c1)

}

// TestSameKey should check that locking and specific key in hash map does not block other lock operations
func TestSameKey(t *testing.T) {
	hashLock := (HashLock{}).New()
	c1 := make(chan string, 2)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		key := hashLock.GetLockKey([]string{"test"})
		hashLock.Lock(key)
		t.Log("Got lock 1")
		defer func() {
			t.Log("Released lock 1")
			c1 <- "result 1"
			hashLock.Unlock(key)
			wg.Done()
		}()
		time.Sleep(5 * time.Second)
	}()
	//make sure the 1st goroutine runs first
	time.Sleep(1 * time.Second)
	wg.Add(1)
	go func() {
		key := hashLock.GetLockKey([]string{"test"})
		hashLock.Lock(key)
		t.Log("Got lock 2")
		defer func() {
			t.Log("Released lock 2")
			c1 <- "result 2"
			hashLock.Unlock(key)
			wg.Done()
		}()
		time.Sleep(1 * time.Second)
	}()
	wg.Wait()
	//get the first channel value, it should be result 1
	if value := <-c1 ; value != "result 1" {
		t.Error("2 returned first")
	}
	//get the latest channel value it should be result 2
	if value := <-c1 ; value != "result 2" {
		t.Error("1 returned last")
	}
	close(c1)

}


// TestMultipleRandomOperations should check that locking and specific key in hash map does not block other lock operations
func TestMultipleRandomOperations(t *testing.T) {
	hashLock := (HashLock{}).New()
	iterations := 1000000
	c1 := make(chan string, iterations)
	wg := &sync.WaitGroup{}
	for i:=0;i<iterations;i++ {
		wg.Add(1)
		go func() {
			key := hashLock.GetLockKey([]string{"test"})
			hashLock.Lock(key)
			t.Log("Got lock 1")
			defer func() {
				t.Log("Released lock 1")
				c1 <- "result 1"
				hashLock.Unlock(key)
				wg.Done()
			}()
			time.Sleep(5 * time.Second)
		}()

	}
	wg.Wait()
	//get the first channel value, it should be result 1
	if value := <-c1 ; value != "result 1" {
		t.Error("2 returned first")
	}
	//get the latest channel value it should be result 2
	if value := <-c1 ; value != "result 2" {
		t.Error("1 returned last")
	}
	close(c1)

}