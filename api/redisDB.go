package api

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisDB struct {
	client *redis.Client
}

//NewRedisDB create redis client
func NewRedisDB(address, password, dbname string) (kv KV, err error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	kv = &redisDB{
		client: rdb,
	}
	return
}

// Get the item with the provided key.
func (m *redisDB) Get(key string) (val []byte, err error) {
	val, err = m.client.Get(context.Background(), key).Bytes()
	if err != nil {
		panic(err)
	}
	return
}

// Add writes the given item, if no value already exists for its key.
func (m *redisDB) Add(key string, val []byte, expiration time.Duration) (err error) {

	//NX -- Only set the key if it does not already exist.
	_, err = m.client.SetNX(context.Background(), key, string(val), expiration).Result()
	if err != nil {
		panic(err)
	}

	return
}

// Set writes the given item, unconditionally.
func (m *redisDB) Set(key string, val []byte, expiration time.Duration) (err error) {

	err = m.client.Set(context.Background(), key, string(val), expiration).Err()
	if err != nil {
		panic(err)
	}

	return
}

// Delete deletes the item with the provided key.
// return nil error if the item didn't already exist in the cache.
func (m *redisDB) Delete(key string) (err error) {
	err = m.client.Del(context.Background(), key).Err()
	if err != nil {
		panic(err)
	}

	return
}
