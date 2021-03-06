package cacher

import (
	"bytes"
	"testing"
	"time"

	engine "github.com/fresh8/go-cache/engine/memory"
)

// TODO: This should be replaced by a mock, not use memory engine

func TestCacher_Get(t *testing.T) {
	e := engine.NewMemoryStore(time.Second * 60)
	cache := NewCacher(e, 5, 5)
	count := 0
	content := []byte("hello")
	regenerate := func() ([]byte, error) {
		count = count + 1
		return content, nil
	}

	data, err := cache.Get("existing", time.Now().Add(1*time.Minute), regenerate)()
	if err != nil {
		t.Fatalf("no error expected, %s given", err)
	}

	if count != 1 {
		t.Fatalf("regenerate function run count should be 1, %d given", count)
	}

	if bytes.Compare(data, content) != 0 {
		t.Fatalf("data expected to be different, %s expected, %s given", content, data)
	}

	data, err = cache.Get("existing", time.Now().Add(1*time.Minute), regenerate)()
	if bytes.Compare(data, content) != 0 {
		t.Fatalf("data expected to be different, %s expected, %s given", content, data)
	}

	if count != 1 {
		t.Fatalf("regenerate function run count should be 1, %d given", count)
	}

	e.Expire("existing")

	newContent := append(content, []byte("-world")...)
	data, err = cache.Get("existing", time.Now().Add(1*time.Minute), func() ([]byte, error) {
		count = count + 1
		return newContent, nil
	})()

	if bytes.Compare(data, newContent) != 0 {
		t.Fatalf("data expected to be different, %s expected, %s given", newContent, data)
	}

	if count != 2 {
		t.Fatalf("regenerate function run count should be 1, %d given", count)
	}
}
