package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisDB struct {
	client *redis.Client
	ctx    context.Context
}

//NewRedisDB create redis client
func NewRedisDB(address, password string) (kv KV, err error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	kv = &redisDB{
		client: rdb,
		ctx:    context.TODO(),
	}
	return
}

func (m *redisDB) WithContext(ctx context.Context) KV {
	m.ctx = ctx
	return m
}

// Get the item with the provided key.
func (m *redisDB) Get(key string) (val []byte, err error) {
	val, err = m.client.Get(m.ctx, key).Bytes()

	return
}

// Add writes the given item, if no value already exists for its key.
func (m *redisDB) Add(key string, val []byte, expiration time.Duration) (err error) {

	//NX -- Only set the key if it does not already exist.
	_, err = m.client.SetNX(m.ctx, key, string(val), expiration).Result()

	return
}

// Set writes the given item, unconditionally.
func (m *redisDB) Set(key string, val []byte, expiration time.Duration) (err error) {

	return m.client.Set(m.ctx, key, string(val), expiration).Err()
}

// Incr increment key
func (m *redisDB) Incr(key string) (err error) {
	return m.client.Incr(m.ctx, key).Err()
}

// Delete deletes the item with the provided key.
func (m *redisDB) Delete(key string) (err error) {
	err = m.client.Del(m.ctx, key).Err()

	return
}
