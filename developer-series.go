package main

import (
	"developer-series/config"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DeveloperSeriesStackProps struct {
	awscdk.StackProps
}

func NewDeveloperSeriesStack(scope constructs.Construct, id string, props *DeveloperSeriesStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create lambda function
	lambdaHandler := awslambda.NewFunction(stack, jsii.String(config.FunctionName), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-" + config.FunctionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(config.MemorySize),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(config.MaxDuration)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String(config.CodePath), nil),
		Handler:      jsii.String(config.Handler),
	})

	// Create API Gateway rest api.
	restApi := awsapigateway.NewRestApi(stack, jsii.String("LambdaRestApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(*stack.StackName() + "-LambdaRestApi"),
		Description: jsii.String("AWS Developer Series REST API"),
	})

	// Add path resources to rest api
	devApiRes := restApi.Root().AddResource(jsii.String("developer-series"), nil)
	devApiRes.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(lambdaHandler, nil), &awsapigateway.MethodOptions{
		ApiKeyRequired: jsii.Bool(false),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewDeveloperSeriesStack(app, config.StackName, &DeveloperSeriesStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	account := os.Getenv("CDK_DEPLOY_ACCOUNT")
	region := os.Getenv("CDK_DEPLOY_REGION")

	if len(account) == 0 || len(region) == 0 {
		account = os.Getenv("CDK_DEFAULT_ACCOUNT")
		region = os.Getenv("CDK_DEFAULT_REGION")
	}

	return &awscdk.Environment{
		Account: jsii.String(account),
		Region:  jsii.String(region),
	}
}
