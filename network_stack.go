package main

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NetworkStack(scope constructs.Construct, id string, props *InfrastructureCdkSampleGoStackProps, networkProps *ec2.VpcProps) cdk.Stack {
	var stackProps cdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}

	networkStack := cdk.NewStack(scope, &id, &stackProps)

	vpc := ec2.NewVpc(networkStack, jsii.String("Vpc"), networkProps)

	cdk.NewCfnOutput(networkStack, jsii.String("VpcId"), &cdk.CfnOutputProps{Value: vpc.VpcId()})

	return networkStack

}
