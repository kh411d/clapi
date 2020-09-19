package db

import (
	"context"
	"time"
)

//KV key value db interface
type KV interface {
	Add(string, []byte, time.Duration) error
	Set(string, []byte, time.Duration) error
	Delete(string) error
	Get(string) ([]byte, error)
	Incr(string) error
	IncrBy(string, int64) error
	WithContext(context.Context) KV
}

//Clap data document structure
type Clap struct {
	URL   string `fauna:"url"`
	Count int64  `fauna:"count"`
}
