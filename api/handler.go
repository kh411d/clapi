package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var db KV
var conf *viper.Viper

func init() {
	var err error
	conf = viper.New()
	conf.AutomaticEnv()

	db, err = NewRedisDB(
		cast.ToString(conf.Get("REDIS_HOST")),
		cast.ToString(conf.Get("REDIS_PASSWORD")),
	)

	if err != nil {
		panic(err)
	}

}

//ReqData request data that comes from post body
type ReqData struct {
	URL string `json:"url"`
}

func handleGetClap(ctx context.Context, db KV, k string) (*events.APIGatewayProxyResponse, error) {
	v, err := db.WithContext(ctx).Get(k)
	if err != nil {
		log.Printf("Error getClap: %v", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to process payload: %v", err),
		}, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(v),
	}, nil

}

func handleAddClap(ctx context.Context, db KV, k string) (*events.APIGatewayProxyResponse, error) {

	if err := db.WithContext(ctx).Incr(k); err != nil {
		log.Printf("Error addClap: %v", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to process payload: %v", err),
		}, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "ok",
	}, nil
}

func httpError(w http.ResponseWriter, httpcode int, err error) {

	res := events.APIGatewayProxyResponse{
		StatusCode: httpcode,
		Body:       fmt.Sprintf("Failed to parse payload: %v", err),
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//ServeHTTP net/http handler
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var res *events.APIGatewayProxyResponse
	var err error

	switch r.Method {
	case "GET":
		url := r.URL.Query().Get("url")
		res, err = handleGetClap(r.Context(), db, url)
		if err != nil {
			httpError(w, 500, err)
			return
		}

	case "POST":
		var data ReqData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			httpError(w, 500, err)
			return
		}

		res, err = handleAddClap(r.Context(), db, data.URL)
		if err != nil {
			httpError(w, 500, err)
			return
		}

	default:
		httpError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//ServeLambda AWS lambda event handler
func ServeLambda(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	switch r.HTTPMethod {
	case "GET":
		url := cast.ToString(r.QueryStringParameters["url"])
		return handleGetClap(context.Background(), db, url)
	case "POST":
		var data ReqData

		if err := json.NewDecoder(strings.NewReader(r.Body)).Decode(&data); err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Failed to parse payload: %v", err),
			}, nil
		}
		return handleAddClap(context.Background(), db, data.URL)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       fmt.Sprintf("Failed to parse payload: %v", "Method not allowed"),
	}, nil
}

//KV key-value storage interface
type KV interface {
	Add(key string, val []byte, expiration time.Duration) error
	Set(key string, val []byte, expiration time.Duration) error
	Delete(key string) error
	Get(key string) ([]byte, error)
	Incr(key string) error
	WithContext(ctx context.Context) KV
}

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
