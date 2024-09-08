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
		EventBridgeEnabled: jsii.Bool(true), // Enable EventBridge integration
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

	// Step 3: Create an EventBridge rule to capture S3 events
	eventRule := awsevents.NewRule(stack, jsii.String("S3EventRule"), &awsevents.RuleProps{
		EventPattern: &awsevents.EventPattern{
			Source:     jsii.Strings("aws.s3"),
			DetailType: jsii.Strings("Object Created"),
			Detail: &map[string]interface{}{
				"bucket": map[string]interface{}{
					"name": jsii.Strings(*bucket.BucketName()), // Bucket name must be an array
				},
				"object": map[string]interface{}{
					"key": map[string]interface{}{
						"prefix": jsii.Strings("order/"), // Prefix must be an array
					},
				},
			},
		},
	})

	// Step 4: Add Lambda function as the target for EventBridge rule
	eventRule.AddTarget(awseventstargets.NewLambdaFunction(processingLambda, nil))

	// Step 5: Add permission for EventBridge to invoke the Lambda function
	processingLambda.AddPermission(jsii.String("AllowEventBridgeInvoke"), &awslambda.Permission{
		Action:    jsii.String("lambda:InvokeFunction"),
		Principal: awsiam.NewServicePrincipal(jsii.String("events.amazonaws.com"), nil),
		SourceArn: eventRule.RuleArn(), // Allow invocation only from this rule
	})

	// Step 6: Add permissions for EventBridge to access S3 events
	bucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:PutObject", "s3:GetObject", "s3:ListBucket"),
		Resources: jsii.Strings(*bucket.BucketArn(), *bucket.BucketArn()+"/*"),
		Principals: &[]awsiam.IPrincipal{
			awsiam.NewServicePrincipal(jsii.String("events.amazonaws.com"), nil),
		},
	}))

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
