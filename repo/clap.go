package repo

import (
	"context"
	"log"

	"github.com/kh411d/clapi/db"
)

//Clap expose Clapper
var Clap Clapper

func init() {
	Clap = &claps{}
}

//Clapper interface
type Clapper interface {
	GetClap(context.Context, db.KV, string) string
	AddClap(context.Context, db.KV, string, int64)
}

type claps struct{}

//GetClap get a clap data
func (claps) GetClap(ctx context.Context, dbconn db.KV, k string) string {
	v, err := dbconn.WithContext(ctx).Get(k)
	if err != nil {
		log.Printf("Error getClap: %v", err)
	}

	if v == nil {
		v = []byte("0")
	}

	//Do not error
	return string(v)
}

//AddClap increament clap data count
func (claps) AddClap(ctx context.Context, dbconn db.KV, k string, v int64) {
	if err := dbconn.WithContext(ctx).IncrBy(k, v); err != nil {
		log.Printf("Error addClap: %v", err)
	}
}
