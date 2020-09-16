package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-redis/redis/v8"
)

func getClaps() map[string]interface{} {
	return map[string]interface{}{
		"StatusCode": 200,
		"Body":       "Get clap",
	}
}

func addClap() map[string]interface{} {
	return map[string]interface{}{
		"StatusCode": 200,
		"Body":       "Adding clap",
	}
}

//HTTPHandler net/http handler
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	var res map[string]interface{}

	switch r.Method {
	case "GET":
		res = getClaps()
	case "POST":
		res = addClap()
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

//EventHandler AWS lambda event handler
func EventHandler(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res map[string]interface{}
	switch r.HTTPMethod {
	case "GET":
		res = getClaps()
	case "PUT":
		res = addClap()
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: res["StatusCode"].(int),
		Body:       res["Body"].(string),
	}, nil
}

//KV key-value storage interface
type KV interface {
	Add(key string, val []byte, expiration time.Duration) error
	Set(key string, val []byte, expiration time.Duration) error
	Delete(key string) error
	Get(key string) ([]byte, error)
}

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
