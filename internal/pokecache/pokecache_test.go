package pokecache

import (
	"testing"
	"time"
)

func TestCreateCache(t *testing.T) {
	cache := NewCache(time.Millisecond)
	if cache.cache == nil {
		t.Error("Cache was not created")
	}
}

func TestAddGetCache(t *testing.T) {
	cache := NewCache(time.Millisecond)

	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "k1",
			value: []byte("v1"),
		},
		{
			key:   "k2",
			value: []byte("v2"),
		},
		{
			key:   "",
			value: []byte("v3"),
		},
	}

	for _, c := range cases {
		cache.Add(c.key, []byte(c.value))
		v, ok := cache.Get(c.key)
		if !ok {
			t.Errorf("%s was not found", c.key)
			continue
		}

		if string(v) != string(c.value) {
			t.Errorf("The value of %s doesn't match. Expected %s, but found %s.", c.key, c.value, v)
			continue
		}
	}
}

func TestReap(t *testing.T) {
	interval := 10 * time.Millisecond
	cache := NewCache(interval)

	k := "key 1"
	cache.Add(k, []byte("value 1"))

	time.Sleep(interval / 2)
	_, ok := cache.Get(k)
	if !ok {
		t.Errorf("%s should not have been reaped", k)
	}

	time.Sleep(interval/2 + time.Millisecond)
	_, ok = cache.Get(k)
	if ok {
		t.Errorf("%s should have been reaped", k)
	}
}

