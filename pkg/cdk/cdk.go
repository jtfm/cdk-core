package cdk

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ScopedCfnOutputProps struct {
	awscdk.CfnOutputProps
	name  *string
	value *string
}

func NewScopedCfnOutput(
	scope constructs.Construct,
	id *string,
	props *ScopedCfnOutputProps) awscdk.CfnOutput {

	var exportName string
	if props.ExportName != nil {
		exportName = *props.ExportName
	} else {
		// Check if the scope is a stack
		if stack, ok := (scope).(awscdk.Stack); ok {
			exportName = fmt.Sprintf("%s-%s", *props.name, *stack.StackName())
		} else {
			exportName = fmt.Sprintf("%s-%s", *props.name, *scope.Node().Id())
		}
	}

	return awscdk.NewCfnOutput(
		scope,
		jsii.String(*props.name),
		&awscdk.CfnOutputProps{
			Value:      props.value,
			ExportName: &exportName,
		})
}
