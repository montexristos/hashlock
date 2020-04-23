# hashlock

## Description

Package hashlock provides a way to lock multiple string indexes using a single hashmap and only allow write access to one thread at a time

Battle-tested in multiple threaded applications with thousand concurrent writes and reads per minute

## Background

Most lock mechanisms in golang use sync.Mutex which works quite well most of the times.

However, if someone tries to acquire an uninitialized Mutex non recoverable error will occur, and checking if the Mutex is initialized may require complex custom logic

That is the reason this library uses boolean channels to handle locking which handle pretty well in terms of stability and performance

## IMPORTANT NOTE

This library is not designed to be a queue, so do not expect a FIFO or any other behavior

The goal is to avoid multiple threads writing to the same resource, while keeping it available for concurrent reads

## Usage

    import "github.com/montexristos/hashlock"

    hashLock := (&HashLock{}).New()

    // lock the value for read/write
    hashLock.Lock("awesomeKey")
    defer hashLock.Unlock("awesomeKey")
    
    //
    go hashLock.Lock("awesomeKey")

    // the goroutine initiated above will wait for the current function to release the lock
