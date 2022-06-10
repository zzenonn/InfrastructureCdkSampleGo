package main

import (
	"fmt"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NetworkStack(scope constructs.Construct, id string, props *InfrastructureCdkSampleGoStackProps, networkProps *ec2.VpcProps) ec2.Vpc {
	var stackProps cdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}

	networkStack := cdk.NewStack(scope, &id, &stackProps)

	vpc := ec2.NewVpc(networkStack, jsii.String("Vpc"), networkProps)

	privateSubnets := vpc.PrivateSubnets()

	isolatedSubnets := vpc.IsolatedSubnets()

	isolatedNacl := ec2.NewNetworkAcl(networkStack, jsii.String("DbNacl"), &ec2.NetworkAclProps{
		Vpc:             vpc,
		SubnetSelection: &ec2.SubnetSelection{Subnets: isolatedSubnets},
	})

	// Ingress rules from private subnet to db subnet
	for index, subnet := range *privateSubnets {
		resourceName := fmt.Sprintf("DbNaclIngress%d", (1+index)*100)
		isolatedNacl.AddEntry(&resourceName, &ec2.CommonNetworkAclEntryOptions{
			RuleNumber: jsii.Number(float64(index * 100)),
			Cidr:       ec2.AclCidr_Ipv4(subnet.Ipv4CidrBlock()),
			Traffic:    ec2.AclTraffic_TcpPort(jsii.Number(5432)), //postgres port
			Direction:  ec2.TrafficDirection_INGRESS,
			RuleAction: ec2.Action_ALLOW,
		})
	}

	for index, subnet := range *privateSubnets {
		resourceName := fmt.Sprintf("DbNaclEgress%d", (1+index)*100)
		isolatedNacl.AddEntry(&resourceName, &ec2.CommonNetworkAclEntryOptions{
			RuleNumber: jsii.Number(float64(index * 100)),
			Cidr:       ec2.AclCidr_Ipv4(subnet.Ipv4CidrBlock()),
			Traffic:    ec2.AclTraffic_TcpPortRange(jsii.Number(1024), jsii.Number(65535)), //Dynamic ports
			Direction:  ec2.TrafficDirection_EGRESS,
			RuleAction: ec2.Action_ALLOW,
		})
	}

	cdk.NewCfnOutput(networkStack, jsii.String("VpcId"), &cdk.CfnOutputProps{Value: vpc.VpcId()})

	return vpc

}
