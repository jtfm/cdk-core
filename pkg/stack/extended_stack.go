package stack

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/iancoleman/strcase"
)

type Build uint8

const (
	Prod Build = iota
	Dev
)

type StackProps struct {
	awscdk.StackProps
	Semver semver.Version
	Build  Build
}

func (s Build) String() string {
	return [...]string{"prod", "dev"}[s]
}

type Stack struct {
	awscdk.Stack
}

func NewStack(
	scope constructs.Construct,
	id *string,
	props *StackProps) Stack {

	if props == nil {
		panic("props are required")
	}

	props.StackName = jsii.String(fmt.Sprintf(
		"%s-%s-%s",
		strcase.ToKebab(*props.StackName),
		props.Semver.String(),
		props.Build.String()))

	stack := awscdk.NewStack(scope, id, &props.StackProps)

	return Stack{stack}
}

func (s Stack) GetName() *string {
	return s.StackName()
}
