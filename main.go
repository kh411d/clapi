package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kh411d/clapi/api"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	v.AutomaticEnv()

	fmt.Print(v.Get("HOME"))

	if v.Get("AWS_LAMBDA_FUNCTION_NAME") == "" {
		//	lambda.Start(handler)
		http.HandleFunc("/clap", api.HttpHandler)
		http.ListenAndServe(":3000", nil)
	} else {
		// Make the handler available for Remote Procedure Call by AWS Lambda
		lambda.Start(api.EventHandler)
	}

}
