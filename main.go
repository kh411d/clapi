package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kh411d/clapi/api"
)

func main() {

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		http.HandleFunc("/", api.ServeHTTP)
		http.ListenAndServe(":3000", nil)
	} else {
		// Make the handler available for Remote Procedure Call by AWS Lambda
		lambda.Start(api.ServeLambda)
	}

}
