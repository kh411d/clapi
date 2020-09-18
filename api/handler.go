package api

import (
	"context"
	"errors"
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
	var err error
	conf = viper.New()
	conf.AutomaticEnv()

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
		panic(err)
	}

}

func validateURL(urlstr string) bool {
	if conf.Get("URL_HOST") == nil {
		return true
	}

	u, err := url.Parse(urlstr)
	if err != nil {
		log.Printf("Error validateURL: %v", err)
		return false
	}

	if u.Hostname() != conf.Get("URL_HOST") {
		return false
	}

	return true
}

func apiGatewayProxyResponse(code int, body string, err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
	}, err
}

//ServeHTTP net/http handler
func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	urlstr := r.URL.Query().Get("url")

	if !validateURL(urlstr) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	switch r.Method {
	case "GET":
		w.Write([]byte(cast.ToString(repo.GetClap(r.Context(), dbconn, urlstr))))
	case "POST":
		w.Write([]byte(cast.ToString(repo.AddClap(r.Context(), dbconn, urlstr))))
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

//ServeLambda AWS lambda event handler
func ServeLambda(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	urlstr := cast.ToString(r.QueryStringParameters["url"])

	if !validateURL(urlstr) {
		return apiGatewayProxyResponse(
			http.StatusForbidden,
			http.StatusText(http.StatusForbidden),
			nil,
		)
	}

	switch r.HTTPMethod {
	case "GET":
		body := repo.GetClap(context.Background(), dbconn, urlstr)
		return apiGatewayProxyResponse(200, cast.ToString(body), nil)
	case "POST":
		body := repo.AddClap(context.Background(), dbconn, urlstr)
		return apiGatewayProxyResponse(200, cast.ToString(body), nil)
	}

	return apiGatewayProxyResponse(
		http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed),
		nil,
	)
}
