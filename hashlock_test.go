package tools

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestHashLock(t *testing.T) {
	l := HashLock{}.New()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go iterate(l, i, t, &wg)
	}
	wg.Wait()
}

func iterate(l HashLock, i int, t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(10-2*i) * time.Millisecond)
	l.Lock("key")
	t.Log(fmt.Sprintf("get Lock %d", i))
	time.Sleep(time.Duration(10-2*i) * time.Millisecond)
	l.Unlock("key")
	t.Log(fmt.Sprintf("unLock %d", i))
}

func TestUnlockOfUnlockedResource(t *testing.T) {
	time.AfterFunc(1*time.Second, func() {
		t.Error("could not get lock")
	})
	l := HashLock{}.New()
	l.Unlock("test")
	l.GetLock("test")
}

func TestMultipleRoutinesSameResource(t *testing.T) {

	f, _ := os.Create("test.txt")
	l := HashLock{}.New()
	wg := (&sync.WaitGroup{})
	for i:=0; i < 10000; i++ {
		wg.Add(1)
		go writeToFile(fmt.Sprintf("goroutine %d\r", i), wg, l)
	}
	wg.Wait()

	f, _ = os.Open("test.txt")
	inputReader := bufio.NewReader(f)
	scanner := bufio.NewScanner(inputReader)
	// Count the lines. they should beb 10001
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%v\r", line)
		count++
	}
	//_ = os.Remove("test.txt")
	if count != 10000 {
		fmt.Printf("Invalid count of lines %d", count)
		t.Fail()
	}
}

func writeToFile(text string, wg *sync.WaitGroup, l HashLock) {
	l.Lock("test")
	fmt.Printf("Locked: %s", text)
	defer func() {
		l.Unlock("test")
		fmt.Printf("Unlocked: %s", text)
		wg.Done()
	}()
	f, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("error opening file", err.Error())
	} else {
		_, _ = fmt.Fprintln(f, text)
		_ = f.Close()
	}
}

func TestMultipleRoutinesSameResourceNoHashLock(t *testing.T) {
	f, _ := os.Create("test.txt")
	wg := (&sync.WaitGroup{})
	for i:=0; i < 10000; i++ {
		wg.Add(1)
		go func(t string) {
			_, _ = fmt.Fprintln(f, t)
			wg.Done()
		}(fmt.Sprintf("goroutine %d\r", i))
	}
	wg.Wait()
	_ = f.Close()

	f, _ = os.Open("test.txt")
	inputReader := bufio.NewReader(f)
	scanner := bufio.NewScanner(inputReader)
	// Count the lines. they should beb 10001
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%v\r", line)
		count++
	}
	//_ = os.Remove("test.txt")
	if count != 10000 {
		fmt.Printf("Invalid count of lines %d", count)
		t.Fail()
	}
}

func TestHashLock_DeleteKey(t *testing.T) {
	l := (HashLock{}).New()
	l.GetLock("test")
	if _, found := l.locks["test"]; !found {
		t.Fail()
	}
	l.DeleteKey("test")
	if _, found := l.locks["test"]; found {
		t.Fail()
	}
}

func TestHashLock_GetLock(t *testing.T) {
	l := (HashLock{}).New()
	l.GetLock("test")
	if _, found := l.locks["test"]; !found {
		t.Fail()
	}
}

func TestHashLock_Lock(t *testing.T) {

}

func TestHashLock_Unlock(t *testing.T) {

}