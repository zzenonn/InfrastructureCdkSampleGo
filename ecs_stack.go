package main

import (
	"fmt"
	"log"
	"os"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecrassets "github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	ecspatterns "github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewEcsStack(scope constructs.Construct, id string, props *InfrastructureCdkSampleGoStackProps, ecsProps *ecspatterns.ApplicationLoadBalancedFargateServiceProps,
	vpc ec2.Vpc) cdk.Stack {

	var stackProps cdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}

	ecsStack := cdk.NewStack(scope, &id, &stackProps)

	ecsCluster := ecs.NewCluster(ecsStack, jsii.String("GlobomanticsCluster"), &ecs.ClusterProps{
		Vpc: vpc,
	})

	cwdir, err := os.Getwd()

	if err != nil {
		log.Fatalf("No such directory, %v", err)
	}

	imageAsset := ecrassets.NewDockerImageAsset(ecsStack, jsii.String("Globomantics-Landing-Page"), &ecrassets.DockerImageAssetProps{
		Directory: jsii.String(fmt.Sprintf("%s/globomantics-container-app", cwdir)),
	})

	image := ecs.ContainerImage_FromDockerImageAsset(imageAsset)

	ecsProps.Cluster = ecsCluster
	ecsProps.TaskImageOptions.Image = image

	ecsPattern := ecspatterns.NewApplicationLoadBalancedFargateService(ecsStack, jsii.String("GlobomanticsFargate"), ecsProps)

	cdk.NewCfnOutput(ecsStack, jsii.String("LoadBalancerUrl"), &cdk.CfnOutputProps{Value: ecsPattern.LoadBalancer().LoadBalancerDnsName()})

	return ecsStack

}
