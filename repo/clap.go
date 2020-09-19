package repo

import (
	"context"
	"log"

	"github.com/kh411d/clapi/db"
)

//GetClap get a clap data
func GetClap(ctx context.Context, dbconn db.KV, k string) string {
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
func AddClap(ctx context.Context, dbconn db.KV, k string, v int64) bool {

	if err := dbconn.WithContext(ctx).IncrBy(k, v); err != nil {
		log.Printf("Error addClap: %v", err)
		//Do not error
		return false
	}

	return true
}
