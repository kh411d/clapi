package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockDB "github.com/kh411d/clapi/db/mocks"
	"github.com/kh411d/clapi/repo"
	mockRepo "github.com/kh411d/clapi/repo/mocks"
)

func TestInit(t *testing.T) {
	var isPanic bool
	defer func() {
		if r := recover(); r != nil {
			isPanic = true
		}
	}()
	os.Setenv("FAUNADB_SECRET_KEY", "")
	os.Setenv("REDIS_HOST", "")

	setEnv()
	setDB()

	if !isPanic {
		t.Errorf("Should be panic")
	}

	os.Setenv("FAUNADB_SECRET_KEY", "testing")
	os.Setenv("REDIS_HOST", "")

	setEnv()
	setDB()

	if dbconn == nil {
		t.Errorf("DBCon should be created")
	}

	os.Setenv("FAUNADB_SECRET_KEY", "")
	os.Setenv("REDIS_HOST", "testing")

	setEnv()
	setDB()

	if dbconn == nil {
		t.Errorf("DBCon should be created")
	}
}

func TestGetUrl(t *testing.T) {
	var x string
	x = getURL("://notaurl")
	assert.Equal(t, x, "", "url should be empty string")

	x = getURL("https://khal.web.id")
	u, _ := url.Parse(x)
	assert.Equal(t, u.Scheme, "http", "Scheme should be a HTTP")

	os.Setenv("URL_HOST", "khal.web.id")

	x = getURL("https://khal.web.id")
	assert.Equal(t, x, "http://khal.web.id", "url missmatch")

	x = getURL("https://NOT.khal.web.id")
	assert.Equal(t, x, "", "Url should be empty string")
}

func TestServe(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "9")
	}))
	defer ts.Close()

	mockClapper := &mockRepo.Clapper{}
	mockClapper.On("GetClap", mock.Anything, mock.Anything, mock.Anything).Return("9")
	mockClapper.On("AddClap", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	mockDBConn := &mockDB.KV{}
	mockDBConn.On("WithContext", mock.Anything).Return(mockDBConn)
	mockDBConn.On("Get", mock.Anything).Return([]byte("9"), nil)
	mockDBConn.On("IncrBy", mock.Anything, mock.Anything).Return(nil)

	repo.Clap = mockClapper

	req := httptest.NewRequest("GET", ts.URL+"?url=http://khal.web.id", nil)
	w := httptest.NewRecorder()

	ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200, "status should be ok")
	assert.Equal(t, string(body), "9", "content missmatch")

	//No url
	req = httptest.NewRequest("GET", ts.URL, nil)
	w = httptest.NewRecorder()

	ServeHTTP(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, http.StatusForbidden, "should be forbid, no url query string")

	preq := httptest.NewRequest("POST", ts.URL+"?url=http://khal.web.id", strings.NewReader("9"))
	pw := httptest.NewRecorder()

	ServeHTTP(pw, preq)

	resp = pw.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200, "status should be ok")

	//Serve lambda
	r := events.APIGatewayProxyRequest{
		Body: "9",
		QueryStringParameters: map[string]string{
			"url": "http://khal.web.id",
		},
		HTTPMethod: "GET",
	}

	apiGWResp, _ := ServeLambda(r)
	assert.Equal(t, apiGWResp.StatusCode, 200, "should always 200")
	assert.Equal(t, apiGWResp.Body, "9", "content missmatch")

	r = events.APIGatewayProxyRequest{
		Body:       "9",
		HTTPMethod: "GET",
	}

	apiGWResp, _ = ServeLambda(r)
	assert.Equal(t, apiGWResp.StatusCode, http.StatusForbidden, "should be forbid, no url query string")

	r = events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"url": "http://khal.web.id",
		},
		HTTPMethod: "POST",
	}

	apiGWResp, _ = ServeLambda(r)
	assert.Equal(t, apiGWResp.StatusCode, 200, "should always 200")

}
