package api

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

//HTTPHandler net/http handler
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	var res map[string]interface{}

	switch r.Method {
	case "GET":
		res = getClaps()
	case "PUT":
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
func EventHandler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, World",
	}, nil
}
