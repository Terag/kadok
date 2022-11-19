package cache

import (
	"bytes"
	"testing"
	"time"
)

var TestValue = []byte("{\"key\":\"some json structure\"}")
var TestKey = "test_key"

func TestCreateMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	if len(c.cachedBytes) != 0 {
		t.Errorf("Cached string created not empty")
	}
}

func TestPutMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	c.Put("test", TestValue, time.Duration(60*time.Second))
	if v, ok := c.cachedBytes[TestKey]; ok {
		if !bytes.Equal(v.value, TestValue) {
			t.Errorf("Unexpected value found in cache, expected \"test\" and got \"%s\"", string(v.value))
		}
	}
}

func TestGetFoundMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	c.Put(TestKey, TestValue, time.Duration(60*time.Second))
	if v, ok, _ := c.Get(TestKey); ok {
		if !bytes.Equal(v, TestValue) {
			t.Errorf("Unexpected value found in cache, expected \"test\" and got \"%s\"", string(v))
		}
	} else {
		t.Errorf("Unexpected value found in cache, expected \"test\" and got nothing")
	}
}

func TestGetNotFoundMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	if v, ok, _ := c.Get(TestKey); ok {
		t.Errorf("Unexpected value found in cache, expected nothing and got \"%s\"", string(v))
	}
}

func TestDeleteMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	c.Put(TestKey, TestValue, time.Duration(60*time.Second))
	if _, ok, _ := c.Get(TestKey); ok {
		c.Delete(TestKey)
		if v, ok, _ := c.Get(TestKey); ok {
			t.Errorf("Unexpected value found in cache, expected nothing and got \"%s\"", string(v))
		}
	} else {
		t.Errorf("Unexpected value found in cache, expected \"test\" and got nothing")
	}
}

func TestExistsMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	c.Put(TestKey, TestValue, time.Duration(60*time.Second))
	if ok, _ := c.Exists(TestKey); !ok {
		t.Errorf("Unexpected cache not found, expected to exist and got nothing")
	}
}

func TestNotExistsMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	if ok, _ := c.Exists(TestKey); ok {
		t.Errorf("Unexpected cache found, expected to not exist and got something")
	}
}

func TestExpiredMemoryCache(t *testing.T) {
	c := NewMemoryCache(time.Duration(1 * time.Second))
	defer c.stop()

	c.Put(TestKey, TestValue, 1*time.Second)
	time.Sleep(2 * time.Second)

	if v, ok, _ := c.Get(TestKey); ok {
		t.Errorf("Unexpected value found in cache, expected nothing and got \"%s\"", string(v))
	}
}
