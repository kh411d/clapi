package api

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/kh411d/clapi/db"
	"github.com/kh411d/clapi/repo"
)

var dbconn db.KV
var conf *viper.Viper

func init() {
	setEnv()
	setDB()
}

func setEnv() {
	conf = viper.New()
	conf.AutomaticEnv()
}

func setDB() {
	var err error

	if conf.Get("FAUNADB_SECRET_KEY") != nil {
		dbconn, err = db.NewFaunaDB(
			cast.ToString(conf.Get("FAUNADB_SECRET_KEY")),
			"claps",
			"url_idx",
		)
	} else if conf.Get("REDIS_HOST") != nil {
		//use redis instead
		dbconn, err = db.NewRedisDB(
			cast.ToString(conf.Get("REDIS_HOST")),
			cast.ToString(conf.Get("REDIS_PASSWORD")),
		)
	} else {
		err = errors.New("DB not found")
	}

	if err != nil {
		log.Printf("No DB Env found")
		panic(err)
	}
}

func getURL(urlstr string) string {

	u, err := url.ParseRequestURI(urlstr)
	if err != nil {
		log.Printf("Error validateURL: %v", err)
		return ""
	}

	u.Scheme = "http"

	if conf.Get("URL_HOST") != nil {
		if u.Hostname() != conf.Get("URL_HOST") {
			return ""
		}
	}

	return u.String()
}

func apiGatewayProxyResponse(code int, body string, err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
	}, err
}

//ServeHTTP net/http handler
func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

	urlstr := getURL(r.URL.Query().Get("url"))

	if urlstr == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	switch r.Method {
	case "GET":
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(repo.Clap.GetClap(r.Context(), dbconn, urlstr)))
	case "POST":
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		repo.Clap.AddClap(r.Context(), dbconn, urlstr, cast.ToInt64(string(rBody)))
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte{})
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

//ServeLambda AWS lambda event handler
func ServeLambda(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	urlstr := getURL(cast.ToString(r.QueryStringParameters["url"]))

	if urlstr == "" {
		return apiGatewayProxyResponse(
			http.StatusForbidden,
			http.StatusText(http.StatusForbidden),
			nil,
		)
	}

	switch r.HTTPMethod {
	case "GET":
		body := repo.Clap.GetClap(context.Background(), dbconn, urlstr)
		return apiGatewayProxyResponse(200, body, nil)
	case "POST":
		repo.Clap.AddClap(context.Background(), dbconn, urlstr, cast.ToInt64(r.Body))
		return apiGatewayProxyResponse(200, "", nil)
	}

	return apiGatewayProxyResponse(
		http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed),
		nil,
	)
}
