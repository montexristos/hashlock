package hashlock

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

// check in stdout the behavior of iterating with goroutines with and without hashlock
func TestHashLock(t *testing.T) {
	l := (&HashLock{}).New(0)
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go iterateNoLock(i, t, &wg)
	}
	wg.Wait()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go iterate(l, i, t, &wg)
	}
	wg.Wait()
}

func iterate(l *HashLock, i int, t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(10-2*i) * time.Millisecond)
	l.Lock("key")
	t.Log(fmt.Sprintf("hashlock before %d", i))
	time.Sleep(time.Duration(10-2*i) * time.Millisecond)
	l.Unlock("key")
	t.Log(fmt.Sprintf("haslock after %d", i))
}

func iterateNoLock(i int, t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(10-2*i) * time.Millisecond)
	t.Log(fmt.Sprintf("noLock before %d", i))
	time.Sleep(time.Duration(10-2*i) * time.Millisecond)
	t.Log(fmt.Sprintf("noLock after %d", i))
}

func TestUnlockOfUnlockedResource(t *testing.T) {
	l := (&HashLock{}).New(0)
	l.Lock("test")
	l.Unlock("test")
	l.Unlock("test")
}

func TestUnlockOfUnitializedResource(t *testing.T) {
	l := (&HashLock{}).New(0)
	l.Unlock("test")
}

// Test writing to a single file from multiple goroutines
func TestMultipleRoutinesSameResource(t *testing.T) {
	f, _ := os.Create("test1.txt")
	l := (&HashLock{}).New(0)
	wg := (&sync.WaitGroup{})
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go writeToFile(f, fmt.Sprintf("goroutine %d\r", i), wg, l)
	}
	wg.Wait()
	_ = f.Close()

	f, _ = os.Open("test1.txt")
	inputReader := bufio.NewReader(f)
	scanner := bufio.NewScanner(inputReader)
	// Count the lines. they should be 100
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		t.Logf("%v\r", line)
		count++
	}
	//_ = os.Remove("test.txt")
	if count != 100 {
		t.Errorf("Invalid count of lines %d", count)
		t.Fail()
	}

}

func writeToFile(f *os.File, t string, wg *sync.WaitGroup, l *HashLock) {
	l.Lock("test")
	defer func() {
		l.Unlock("test")
		wg.Done()
	}()
	_, _ = fmt.Fprintln(f, t)
	_ = f.Sync()
}

// Test deleting a key from hash lock map
func TestHashLock_DeleteKey(t *testing.T) {
	l := (&HashLock{}).New(0)
	l.Lock("test")
	if _, found := l.locks["test"]; !found {
		t.Errorf("test lock not found")
		t.Fail()
	}
	l.DeleteKey("test")
	if _, found := l.locks["test"]; found {
		t.Errorf("test lock found")
		t.Fail()
	}
}

func TestHashLock_GetLock(t *testing.T) {
	l := (&HashLock{}).New(0)
	l.Lock("test")
	if _, found := l.locks["test"]; !found {
		t.Fail()
	}
}

// Lock should unlock automatically after 1 second
func TestHashLock_WaitUnlock(t *testing.T) {
	// set timeout to 3 sec
	l := (&HashLock{}).New(2 * time.Second)
	startTime := time.Now()
	l.Lock("test")
	// do not unlock but wait for timeout
	l.Lock("test")
	endTime := time.Now()
	if endTime.Sub(startTime) < 2*time.Second {
		t.Fail()
	}
}

// Test GetLockKey method
func TestHashLock_GetLockKey(t *testing.T) {
	l := (&HashLock{}).New(0)
	params := make([]string, 0)
	params = append(params, "test")
	key := l.GetLockKey(params)
	if key != "test" {
		t.Errorf("error creating key string %s", params)
		t.Fail()
	}
	params = append(params, "1")
	key = l.GetLockKey(params)
	if key != "test-1" {
		t.Errorf("error creating key string %s", params)
		t.Fail()
	}
	params = append(params, "1")
	key = l.GetLockKey(params)
	if key != "test-1-1" {
		t.Errorf("error creating key string %s", params)
		t.Fail()
	}
}

//Test Unlock that assumes locking an unlocked resource is instant
func TestHashLock_Unlock(t *testing.T) {
	l := (&HashLock{}).New(5 * time.Second)
	l.Lock("test")
	l.Unlock("test")
	startTime := time.Now()
	// lock should be instant
	l.Lock("test")
	endTime := time.Now()
	if endTime.Sub(startTime) > 1*time.Millisecond {
		t.Fail()
	}
}

//Test Unlock that assumes locking an unlocked resource is instant
func TestHashLock_Empty(t *testing.T) {
	l := (&HashLock{}).New(10 * time.Second)
	l.Lock("test")
	l.Lock("test1")
	_ = l.Empty()
	if len(l.locks) == 0 {
		t.Errorf("hashlock emptied while locked")
		t.Fail()
	}
	l.Unlock("test")
	l.Unlock("test1")
	_ = l.Empty()
	if len(l.locks) > 0 {
		t.Errorf("hashlock not emptied while locked")
		t.Fail()
	}
}
