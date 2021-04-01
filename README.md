# hashlock

[![codecov](https://codecov.io/gh/montexristos/hashlock/branch/master/graph/badge.svg?token=6AU1RTYQZX)](https://codecov.io/gh/montexristos/hashlock)
[![Go Reference](https://pkg.go.dev/badge/github.com/montexristos/hashlock.svg)](https://pkg.go.dev/github.com/montexristos/hashlock)

## Description

Package hashlock provides a way to lock multiple string indexes and only allow write access to a resource to one thread at a time

Battle-tested in multi-threaded applications with thousands concurrent writes and reads per minute

## Background

Most lock mechanisms in golang use sync.Mutex which works quite well most of the time.

However, if someone tries to acquire an uninitialized Mutex, a non-recoverable error will occur. Checking if the Mutex is initialized may require complex custom logic

That is the reason this library uses boolean channels to handle locking which handle pretty well in terms of stability and performance

## IMPORTANT NOTE

This library is not designed to be a queue, so do not expect a FIFO or any other behavior

The goal is to avoid multiple threads writing to the same resource, while keeping it available for concurrent reads

## Usage

    import "github.com/montexristos/hashlock"

    hashLock := (&HashLock{}).New(1 * time.Second)

    // lock the value for read/write
    hashLock.Lock("awesomeKey")
    defer hashLock.Unlock("awesomeKey")
    
    //
    go hashLock.Lock("awesomeKey")

    // the goroutine initiated above will wait for the current function to release the lock
