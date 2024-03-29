package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	chiproxy "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
)

// RunHandler runs the handler function in an AWS Lambda environment if it detects one and defaults to the local environment if not. The event is ignored if the function is running in a lambda environment.
func RunHandler(
	handler func(
		context.Context,
		interface{}) (interface{}, error),
	event interface{}) error {

	if event == nil {
		return errors.New("event is nil")
	}

	_, err := json.Marshal(event)
	if err != nil {
		return errors.New("failed to marshal event")
	}

	if IsInLambda() {
		lambda.Start(handler)
	} else {
		ctx := context.Background()
		handler(ctx, event)
	}
	return nil
}

// A function that runs an http handler in a lambda environment or a local server
// depending on whether the function is running in a lambda environment or not.
// The addr argument is ignored if the function is running in a lambda environment.
func SwitchingHttpHandler(addr string, handler http.Handler) error {
	if IsInLambda() {
		lambda.Start(handler)
	} else {
		validatePort(addr)
		err := http.ListenAndServe(addr, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

// A function that runs a router in a lambda environment or a local server
// depending on whether the function is running in a lambda environment or not.
// The addr argument is ignored if the function is running in a lambda environment.
func SwitchingRouter(addr string, router *chi.Mux) error {
	if IsInLambda() {
		chiApiProxy := chiproxy.New(router)
		lambda.Start(
			func(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
				proxyRequest, err := chiApiProxy.ProxyWithContext(ctx, req)
				return &proxyRequest, err
			})
	} else {
		errs := validatePort(addr)
		for _, err := range errs {
			return err
		}
		err := http.ListenAndServe(addr, router)
		if err != nil {
			return err
		}
	}
	return nil
}

// Returns true if the function is running in a lambda environment.
func IsInLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func validatePort(port string) []error {
	var errs []error = []error{}
	if port[0:1] != ":" {
		errs = append(errs,
			fmt.Errorf("port %s must start with a colon", port))
	}
	portInt, err := strconv.Atoi(port[1:])
	if err != nil {
		errs = append(errs,
			fmt.Errorf("port %s must be a number", port))
	}
	if portInt < 1 || portInt > 65535 {
		errs = append(errs,
			fmt.Errorf("port %d must be between 1 and 65535", portInt))
	}
	return errs
}

// Creates a Lambda function using Go that will run on ARM64 architecture.
// Other parameters can be supplied via dockerImageFunctionProps. Test Change.
func CreateArmGoDockerLambda(
	stack constructs.Construct,
	name *string,
	codePath *string,
	dockerImageFunctionProps *awslambda.DockerImageFunctionProps) *awslambda.DockerImageFunction {

	// Create a new DockerImageFunctionProps if one doesn't exist
	if dockerImageFunctionProps == nil {
		dockerImageFunctionProps = &awslambda.DockerImageFunctionProps{}
	}

	dockerImageFunctionProps.FunctionName = name
	dockerImageFunctionProps.Code = awslambda.DockerImageCode_FromImageAsset(
		codePath,
		&awslambda.AssetImageCodeProps{},
	)

	dockerImageFunctionProps.Architecture = awslambda.Architecture_ARM_64()

	// Create a map for environment variables if one doesn't exist
	if dockerImageFunctionProps.Environment == nil {
		dockerImageFunctionProps.Environment = &map[string]*string{}
	}
	env := *dockerImageFunctionProps.Environment

	// Add environment variables necessary for running Go under ARM64
	env["CGO_ENABLED"] = jsii.String("0")
	env["GOOS"] = jsii.String("linux")
	env["GOARCH"] = jsii.String("arm64")

	dockerImageFunctionProps.Environment = &env

	dockerImageFunction := awslambda.NewDockerImageFunction(
		stack,
		name,
		dockerImageFunctionProps,
	)

	return &dockerImageFunction
}
