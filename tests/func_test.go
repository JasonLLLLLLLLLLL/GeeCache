package main

import (
	"GeeCache"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f GeeCache.Getter = GeeCache.GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}