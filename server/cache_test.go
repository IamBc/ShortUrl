package main

import (
	"testing"
	"time"
	)

func TestAdd(t *testing.T) {
    cache := NewCache(10)
    cache.Add(`key`, `value`)
}

func TestGetVal(t *testing.T){
    cache := NewCache(10)

    cache.Add(`key`, `value`)
    val := cache.GetVal(`key`)

    if val != `value` {
	t.Fail()
    }
}


func TestExpire(t *testing.T) {
    cache := NewCache(1)
    cache.Add(`key`, `value`)
    time.Sleep(2000 * time.Millisecond)

    val := cache.GetVal(`key`)

    if val == `value` {
	t.Fail()
    }
}
