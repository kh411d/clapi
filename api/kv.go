package api

import "time"

//KV key-value intervace storage
type KV interface {
	Add(key string, val []byte, expiration time.Duration) error
	Set(key string, val []byte, expiration time.Duration) error
	Delete(key string) error
	Get(key string) ([]byte, error)
}
