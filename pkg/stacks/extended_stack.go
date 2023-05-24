package stacks

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/iancoleman/strcase"
)

// ExtendedStack is a stack that adds some additional functionality to the standard stack
type ExtendedStackProps struct {
	awscdk.StackProps
	Semver semver.Version
}

type ExtendedStack struct {
	awscdk.Stack
}

func NewExtendedStack(
	scope constructs.Construct,
	id *string,
	props *ExtendedStackProps) ExtendedStack {

	if props == nil {
		panic("props are required")
	}

	props.StackName = jsii.String(fmt.Sprintf(
		"%s-%s",
		strcase.ToKebab(*props.StackName),
		props.semver.String()))

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	return ExtendedStack{stack}
}

func (s ExtendedStack) GetName() *string {
	return s.StackName()
}
