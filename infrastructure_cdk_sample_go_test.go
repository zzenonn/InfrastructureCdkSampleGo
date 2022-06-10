package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in infrastructure_cdk_sample_go_test.go
func TestNetworkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	vpc := NetworkStack(app, "NetworkStack", &InfrastructureCdkSampleGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	}, CreateVpcProps())

	// THEN
	template := assertions.Template_FromStack(vpc.Stack())

	template.ResourceCountIs(jsii.String("AWS::EC2::VPC"), jsii.Number(1))
}
