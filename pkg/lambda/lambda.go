package lambda

import (
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

// A function that runs a router in a lambda environment or a local server
// depending on whether the function is running in a lambda environment or not.
// The port argument is ignored if the function is running in a lambda environment.
func SwitchingListenAndServe(router func() http.Handler, port string) {
	if IsInLambda() {
		lambda.Start(router())
	} else {
		validatePort(port)
		err := http.ListenAndServe(port, router())
		if err != nil {
			panic(err)
		}
	}
}

// Returns true if the function is running in a lambda environment.
func IsInLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func validatePort(port string) {
	if port[0:1] != ":" {
		panic("Port must start with a colon")
	}
	portInt, err := strconv.Atoi(port[1:])
	if err != nil {
		panic("Port must be a number")
	}
	if portInt < 1 || portInt > 65535 {
		panic("Port must be between 1 and 65535")
	}
}
