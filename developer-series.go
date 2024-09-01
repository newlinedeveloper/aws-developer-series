package main

import (
	"developer-series/config"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
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
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// Create SNS Topic
	topic := awssns.NewTopic(stack, jsii.String("ProcessingTopic"), &awssns.TopicProps{
		TopicName: jsii.String(*stack.StackName() + "-ProcessingTopic"),
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

	// Subscribe Lambda to SNS Topic
	topic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(processingLambda, &awssnssubscriptions.LambdaSubscriptionProps{
		FilterPolicy: &map[string]awssns.SubscriptionFilter{
			"prefix": awssns.SubscriptionFilter_StringFilter(&awssns.StringConditions{
				Allowlist: &[]*string{jsii.String("orders/")},
			}),
		},
	}))

	bucket.AddEventNotification(
		awss3.EventType_OBJECT_CREATED,
		awss3notifications.NewSnsDestination(topic),
		&awss3.NotificationKeyFilter{
			Prefix: jsii.String("orders/"), // Trigger only for files under 'orders' folder
		},
	)

	// Adding the necessary S3 bucket policy for SNS publishing
	bucket.AddToResourcePolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("s3:PutObject"),
			Resources: jsii.Strings(*bucket.BucketArn() + "/*"),
			Principals: &[]awsiam.IPrincipal{
				awsiam.NewServicePrincipal(jsii.String("sns.amazonaws.com"), nil),
			},
		}),
	)

	// Grant the Lambda function permissions to publish to SNS topic
	topic.GrantPublish(processingLambda)

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
