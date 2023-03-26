package ec2

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/iancoleman/strcase"
)

func exportSubnetIds(
	stack constructs.Construct,
	vpc awsec2.IVpc,
	subnetSelection *awsec2.SubnetSelection) {

	subnetIds := vpc.SelectSubnets(subnetSelection).SubnetIds

	for i, subnetId := range *subnetIds {
		exportName := jsii.String(
			strcase.ToKebab(
				fmt.Sprintf(
					"%s-subnet-id-%d",
					subnetSelection.SubnetType,
					i,
				),
			),
		)
		awscdk.NewCfnOutput(
			stack,
			exportName,
			&awscdk.CfnOutputProps{
				Value:      subnetId,
				ExportName: exportName,
			},
		)
	}
}
