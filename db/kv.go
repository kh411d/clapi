package db

import (
	"context"
	"time"
)

//KV key value db interface
type KV interface {
	Add(key string, val []byte, expiration time.Duration) error
	Set(key string, val []byte, expiration time.Duration) error
	Delete(key string) error
	Get(key string) ([]byte, error)
	Incr(key string) error
	WithContext(ctx context.Context) KV
}

//Clap data document structure
type Clap struct {
	URL   string `fauna:"url"`
	Count int    `fauna:"count"`
}
