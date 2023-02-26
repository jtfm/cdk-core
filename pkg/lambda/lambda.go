package lambda

import (
	"os"
	"net/http"
	"github.com/aws/aws-lambda-go/lambda"
)


func SwitchingListenAndServe(router func() http.Handler) {
	if isInLambda() {
		lambda.Start(router())
	} else {
		http.ListenAndServe(":8080", router())
	}
}

func isInLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}