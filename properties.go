package main

import (
	autoscaling "github.com/aws/aws-cdk-go/awscdk/v2/awsautoscaling"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	jsii "github.com/aws/jsii-runtime-go"
)

var userDataScript string = `#!/bin/bash
# Install Apache Web Server and PHP
yum install -y httpd git
# Download Lab files
git clone https://github.com/ps-interactive/lab_aws_implement-auto-scaling-amazon-ecs
mv lab_aws_implement-auto-scaling-amazon-ecs/webapp/* /var/www/html/
# Turn on web server
chkconfig httpd on
service httpd start`

func CreateAsgProps() *autoscaling.AutoScalingGroupProps {
	userData := ec2.MultipartUserData_Custom(jsii.String(userDataScript))
	asgProps := autoscaling.AutoScalingGroupProps{
		Vpc: nil, //to be assigned later
		VpcSubnets: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
		},
		InstanceType: ec2.NewInstanceType(jsii.String("t3.micro")),
		MachineImage: nil, //to be assigned later
		// DesiredCapacity: jsii.Number(2),
		MinCapacity: jsii.Number(2),
		MaxCapacity: jsii.Number(5),
		UserData:    userData,
	}

	return &asgProps
}

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
