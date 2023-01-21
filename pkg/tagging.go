package tagging

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func SetName(taggable awscdk.ITaggable, name string) {
	taggable.Tags().SetTag(
		jsii.String("Name"), 
		jsii.String(name), 
		jsii.Number(100), 
		jsii.Bool(true))
}