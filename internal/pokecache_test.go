package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(1 * time.Second)
	
	cache.Add("key1", []byte("value1"))
	cache.Add("key2", []byte("value2"))

	if val, exists := cache.Get("key1"); !exists || string(val) != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val)
	}
	if val, exists := cache.Get("key2"); !exists || string(val) != "value2" {
		t.Errorf("Expected 'value2', got '%s'", val)
	}

	time.Sleep(2 * time.Second)

	if _, exists := cache.Get("key1"); exists {
		t.Error("Expected 'key1' to be expired, but it still exists")
	}
	if _, exists := cache.Get("key2"); exists {
		t.Error("Expected 'key2' to be expired, but it still exists")
	}
}

func TestAdd(t *testing.T) {
	cache := NewCache(1 * time.Second)
	cache.Add("key", []byte("value"))
	if val, exists := cache.Get("key"); !exists || string(val) != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}
}

func TestGet(t *testing.T) {
	cache := NewCache(1 * time.Second)
	cache.Add("key", []byte("value"))
	if val, exists := cache.Get("key"); !exists || string(val) != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}
	if _, exists := cache.Get("nonexistent"); exists {
		t.Error("Expected 'nonexistent' to not exist, but it does")
	}
}