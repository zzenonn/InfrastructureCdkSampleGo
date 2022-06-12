package main

import (
	"context"
	"log"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	autoscaling "github.com/aws/aws-cdk-go/awscdk/v2/awsautoscaling"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	elb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ec2sdk "github.com/aws/aws-sdk-go-v2/service/ec2"
)

var amazonLinuxAmi = ec2.NewAmazonLinuxImage(&ec2.AmazonLinuxImageProps{
	Generation:     ec2.AmazonLinuxGeneration_AMAZON_LINUX_2,
	Edition:        ec2.AmazonLinuxEdition_STANDARD,
	Virtualization: ec2.AmazonLinuxVirt_HVM,
	Storage:        ec2.AmazonLinuxStorage_GENERAL_PURPOSE,
})

func NewInstanceStack(scope constructs.Construct, id string, props *InfrastructureCdkSampleGoStackProps, asgProps *autoscaling.AutoScalingGroupProps,
	vpc ec2.Vpc, keyName string, useSsh bool, ec2Type string) cdk.Stack {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	var stackProps cdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}

	ec2svc := ec2sdk.NewFromConfig(cfg)

	_, ec2err := ec2svc.DescribeKeyPairs(context.TODO(), &ec2sdk.DescribeKeyPairsInput{KeyNames: []string{keyName}})

	if ec2err != nil {
		log.Printf("Failed to find key pair %s. Creating it instead. %v", keyName, ec2err)

		keyresp, keyerr := ec2svc.CreateKeyPair(context.TODO(), &ec2sdk.CreateKeyPairInput{KeyName: aws.String(keyName)})

		log.Print(*keyresp.KeyMaterial)

		if keyerr != nil {
			log.Fatalf("failed to create key pair, %v", keyerr)
		}
	} else {
		useSsh = false
	}

	instanceStack := cdk.NewStack(scope, &id, &stackProps)

	bastion := ec2.NewBastionHostLinux(instanceStack, jsii.String("BastionHost"), &ec2.BastionHostLinuxProps{
		Vpc:             vpc,
		SubnetSelection: &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PUBLIC},
		InstanceName:    jsii.String("Bastion Host"),
		InstanceType:    ec2.NewInstanceType(&ec2Type),
	})

	if useSsh {
		bastion.Connections().AllowFromAnyIpv4(ec2.NewPort(&ec2.PortProps{
			StringRepresentation: jsii.String("SSH"),
			Protocol:             ec2.Protocol_TCP,
			FromPort:             jsii.Number((22)),
			ToPort:               jsii.Number(22)}),
			jsii.String("Internet access SSH"))

		bastion.Instance().Instance().AddPropertyOverride(jsii.String("KeyName"), keyName)
	}

	ssmPolicy := iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Effect: iam.Effect_ALLOW,
		Resources: &[]*string{
			jsii.String("*"),
		},
		Actions: &[]*string{
			jsii.String("ssmmessages:*"),
			jsii.String("ssm:UpdateInstanceInformation"),
			jsii.String("ec2messages:*"),
		},
	})

	alb := elb.NewApplicationLoadBalancer(instanceStack, jsii.String("WebAlb"), &elb.ApplicationLoadBalancerProps{
		Vpc:              vpc,
		InternetFacing:   jsii.Bool(true),
		LoadBalancerName: jsii.String("WebAlb"),
	})

	httpPort := ec2.NewPort(&ec2.PortProps{
		Protocol:             ec2.Protocol_TCP,
		StringRepresentation: jsii.String("WebPort"),
		FromPort:             jsii.Number(80),
		ToPort:               jsii.Number(80),
	})

	listener := alb.AddListener(jsii.String("Web"), &elb.BaseApplicationListenerProps{
		Port: jsii.Number(80),
		Open: jsii.Bool(true),
	})

	alb.Connections().AllowFromAnyIpv4(httpPort, jsii.String("Allow from internet to port 80 ALB"))

	asgProps.Vpc = vpc
	asgProps.MachineImage = amazonLinuxAmi

	asg := autoscaling.NewAutoScalingGroup(instanceStack, jsii.String("GlobomanticsWeb"), asgProps)

	asg.AddToRolePolicy(ssmPolicy)

	asg.Connections().AllowFrom(alb, httpPort, jsii.String("Allow from ALB to ASG"))

	listener.AddTargets(jsii.String("WebTargetGroup"), &elb.AddApplicationTargetsProps{
		Port: jsii.Number(80),
		Targets: &[]elb.IApplicationLoadBalancerTarget{
			asg,
		},
	})

	cdk.NewCfnOutput(instanceStack, jsii.String("LoadBalancerUrl"), &cdk.CfnOutputProps{Value: alb.LoadBalancerDnsName()})

	return instanceStack

}
