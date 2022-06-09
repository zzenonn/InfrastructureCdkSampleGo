package main

import (
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	jsii "github.com/aws/jsii-runtime-go"
)

func CreateVpcProps() *ec2.VpcProps {

	vpcProps := ec2.VpcProps{
		Cidr:        jsii.String("10.0.0.0/16"),
		MaxAzs:      jsii.Number(6),
		NatGateways: jsii.Number(1),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				Name:       jsii.String("Public"),
				SubnetType: ec2.SubnetType_PUBLIC,
				CidrMask:   jsii.Number(24),
			},
			{
				Name:       jsii.String("Private"),
				SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
				CidrMask:   jsii.Number(24),
			},
			{
				Name:       jsii.String("DB"),
				SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
				CidrMask:   jsii.Number(24),
			},
		},
	}

	return &vpcProps
}
