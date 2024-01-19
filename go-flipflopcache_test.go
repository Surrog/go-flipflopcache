package goflipflopcache

import (
	"testing"
	"time"
)

func TestBasicCache(t *testing.T) {
	cache := NewFlipFlopCache(10 * time.Second)

	cache.Append("test", nil)
	val, ok := cache.Get("test")
	if !ok {
		t.Fatalf("failed to get previously inserted key")
	}

	if val != nil {
		t.Fatalf("unexpected value returned")
	}

	cache.Append("integer", []byte("42"))
	val, ok = cache.Get("integer")
	if !ok {
		t.Fatalf("failed to get previously inserted key")
	}

	if val == nil || string(val) != "42" {
		t.Fatalf("failed to get value %v", val)
	}

	_, ok = cache.Get("toto")
	if ok {
		t.Fatalf("unexpected value inside")
	}
}

func TestFlipFlop(t *testing.T) {
	cache := NewFlipFlopCache(10 * time.Second)

	cache.Append("integer", []byte("1"))
	val, ok := cache.Get("integer")
	if !ok || val == nil || string(val) != "1" {
		t.Fatalf("invalid value inserted %v : %v", ok, string(val))
	}
	cache.Flip()
	val, ok = cache.Get("integer")
	if !ok || val == nil || string(val) != "1" {
		t.Fatalf("invalid value after first flip %v : %v", ok, string(val))
	}
	cache.Flip()
	val, ok = cache.Get("integer")
	if !ok || val == nil || string(val) != "1" {
		t.Fatalf("invalid value after first flip %v : %v", ok, string(val))
	}
	cache.Append("integer", []byte("2"))
	val, ok = cache.Get("integer")
	if !ok || val == nil || string(val) != "2" {
		t.Fatalf("invalid value after first flip %v : %v", ok, string(val))
	}
	cache.Flip()
	val, ok = cache.Get("integer")
	if !ok || val == nil || string(val) != "2" {
		t.Fatalf("invalid value after first flip %v : %v", ok, string(val))
	}
	cache.Flip()
	cache.Flip()
	val, ok = cache.Get("integer")
	if ok {
		t.Fatalf("invalid value after first flip %v : %v", ok, string(val))
	}

	cache.Append("integer", []byte("3"))
	val, ok = cache.Get("integer")
	if !ok || val == nil || string(val) != "3" {
		t.Fatalf("invalid value after first flip %v : %v", ok, string(val))
	}
	time.Sleep(11 * time.Second)
	val, ok = cache.Get("integer")
	if !ok || val == nil || string(val) != "3" {
		t.Fatalf("invalid value after wait flip %v : %v", ok, string(val))
	}
	if time.Since(cache.timeLastFlip) > time.Second {
		t.Fatalf("Flip not happened yet")
	}
}

func TestReset(t *testing.T) {
	cache := NewFlipFlopCache(5 * time.Second)

	cache.Append("test", nil)
	_, ok := cache.Get("test")
	if !ok {
		t.Fatalf("expected value inside")
	}
	time.Sleep(6 * time.Second)
	_, ok = cache.Get("test")
	if !ok {
		t.Fatalf("expected value inside")
	}
	time.Sleep(11 * time.Second)
	_, ok = cache.Get("test")
	if ok {
		t.Fatalf("unexpected value inside")
	}
	cache.Append("test", nil)
	_, ok = cache.Get("test")
	if !ok {
		t.Fatalf("expected value inside")
	}
	cache.Reset()
	_, ok = cache.Get("test")
	if ok {
		t.Fatalf("unexpected value inside")
	}
}

func TestAppendAfterFlip(t *testing.T) {
	cache := NewFlipFlopCache(5 * time.Second)

	time.Sleep(6 * time.Second)
	cache.Append("test", nil)
	_, ok := cache.Get("test")
	if !ok {
		t.Fatalf("expected value inside")
	}
	time.Sleep(11 * time.Second)
	cache.Append("test", nil)
	_, ok = cache.Get("test")
	if !ok {
		t.Fatalf("expected value inside")
	}
}
