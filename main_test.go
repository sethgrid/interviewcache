package interviewcache

import (
	"testing"
	"time"
)

func TestCacheServedCachedItem(t *testing.T) {

	datam := map[string][]byte{
		"foo": []byte("bar"),
		"raz": []byte("qux"),
	}
	fg := &fakeGetter{
		m: datam,
	}
	c := NewCache(fg)

	data, ok := c.Get("foo")
	if !ok {
		t.Fatalf("item not found in cache")
	}
	if got, want := string(data), "bar"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	// update backend
	datam["foo"] = []byte("foo2")
	time.Sleep(20 * time.Millisecond)
	data, ok = c.Get("foo")
	if !ok {
		t.Fatalf("item not found in cache")
	}
	// should still see bar as it is cached
	if got, want := string(data), "bar"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

type fakeGetter struct {
	m map[string][]byte
}

func (f *fakeGetter) Get(key string) ([]byte, error) {
	data := f.m[key]
	return data, nil
}
