package main

import (
	"developer-series/config"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
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

	// Create the S3 bucket
	bucket := awss3.NewBucket(stack, jsii.String("ProcessingBucket"), &awss3.BucketProps{
		RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
		BlockPublicAccess:  awss3.BlockPublicAccess_BLOCK_ALL(),
		EventBridgeEnabled: jsii.Bool(true), // Enable EventBridge integration
	})

	s3ObjectTaggingStatement := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:PutObjectTagging"),
		Resources: jsii.Strings(*bucket.BucketArn() + "/*"),
	})

	// Create lambda function
	processingLambda := awslambda.NewFunction(stack, jsii.String(config.FunctionName), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-" + config.FunctionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(config.MemorySize),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(config.MaxDuration)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String(config.CodePath), nil),
		Handler:      jsii.String(config.Handler),
	})

	processingLambda.AddToRolePolicy(s3ObjectTaggingStatement)
	if processingLambda.Role() != nil {
		processingLambda.Role().AttachInlinePolicy(awsiam.NewPolicy(stack, jsii.String("S3ObjectTaggingPolicy"), &awsiam.PolicyProps{
			Statements: &[]awsiam.PolicyStatement{s3ObjectTaggingStatement},
		}))
	}

	// Create an EventBridge rule for S3 object creation events
	rule := awsevents.NewRule(stack, jsii.String("S3EventRule"), &awsevents.RuleProps{
		EventPattern: &awsevents.EventPattern{
			Source:     jsii.Strings("aws.s3"),
			DetailType: jsii.Strings("Object Created"),
			Detail: &map[string]interface{}{
				"bucket": map[string]interface{}{
					"name": jsii.Strings(*bucket.BucketName()),
				},
			},
		},
	})

	// Add Lambda function as target for the EventBridge rule
	rule.AddTarget(awseventstargets.NewLambdaFunction(processingLambda, &awseventstargets.LambdaFunctionProps{
		MaxEventAge:   awscdk.Duration_Hours(jsii.Number(2)),
		RetryAttempts: jsii.Number(3),
	}))

	// Output bucket name, Lambda function ARN, and EventBridge rule ARN
	awscdk.NewCfnOutput(stack, jsii.String("S3BucketName"), &awscdk.CfnOutputProps{
		Value: bucket.BucketName(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("LambdaFunctionARN"), &awscdk.CfnOutputProps{
		Value: processingLambda.FunctionArn(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("EventBridgeRuleARN"), &awscdk.CfnOutputProps{
		Value: rule.RuleArn(),
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
