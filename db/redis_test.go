package db

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisGet(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	s.Set("foo", "bar")

	kv, _ := NewRedisDB(s.Addr(), "")

	val, _ := kv.WithContext(context.Background()).Get("foo")
	assert.Equal(t, val, []byte("bar"), "missmatch")

	kv.Add("foo", []byte("zap"), 0)
	val, _ = kv.Get("foo")
	assert.Equal(t, val, []byte("bar"), "zap should not be added")

	kv.Add("boo", []byte("1"), 0)
	val, _ = kv.Get("boo")
	assert.Equal(t, val, []byte("1"), "boo should not be added")

	kv.Set("foo", []byte("dug"), 0)
	val, _ = kv.Get("foo")
	assert.Equal(t, val, []byte("dug"), "dug should be added")

	kv.Incr("boo")
	val, _ = kv.Get("boo")
	assert.Equal(t, val, []byte("2"), "boo should be 2 now")

	kv.IncrBy("boo", 3)
	val, _ = kv.Get("boo")
	assert.Equal(t, val, []byte("5"), "boo should be 5 now")

	kv.Delete("foo")
	val, _ = kv.Get("foo")
	assert.Empty(t, val, "foo should be empty now")
}
