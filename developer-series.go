package main

import (
	"developer-series/config"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
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

	// Define the DynamoDB table
	ordersTable := awsdynamodb.NewTable(stack, jsii.String(config.TableName), &awsdynamodb.TableProps{
		TableName: jsii.String(config.TableName),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String(config.PartitionKey),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String(config.SortKey),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// Create records lambda function
	createOrderHandler := CreateLambdaHandler(stack, config.CreateOrderFunctionName, config.CreateOrderCodePath)

	// Read records lambda function
	readOrderHandler := CreateLambdaHandler(stack, config.ReadOrderFunctionName, config.ReadOrderCodePath)

	// Update records lambda function
	updateOrderHandler := CreateLambdaHandler(stack, config.UpdateOrderFunctionName, config.UpdateOrderCodePath)

	// Delete records lambda function
	deleteOrderHandler := CreateLambdaHandler(stack, config.DeleteFunctionName, config.DeleteOrderCodePath)

	// Create API Gateway rest api.
	restApi := awsapigateway.NewRestApi(stack, jsii.String("LambdaRestApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(*stack.StackName() + "-LambdaRestApi"),
		Description: jsii.String("AWS Developer Series REST API"),
	})

	ordersApi := restApi.Root().AddResource(jsii.String("orders"), nil)

	// Add path resources to rest api
	createOrderApiRes := ordersApi.AddResource(jsii.String("create"), nil)
	createOrderApiRes.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(createOrderHandler, nil), nil)

	readOrderApiRes := ordersApi.AddResource(jsii.String("read"), nil)
	readOrderApiRes.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(readOrderHandler, nil), nil)

	updateOrderApiRes := ordersApi.AddResource(jsii.String("update"), nil)
	updateOrderApiRes.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(updateOrderHandler, nil), nil)

	deleteOrderApiRes := ordersApi.AddResource(jsii.String("delete"), nil)
	deleteOrderApiRes.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(deleteOrderHandler, nil), nil)

	// Grant the Lambda function permissions to perform CRUD operations on the DynamoDB table
	ordersTable.GrantWriteData(createOrderHandler)
	ordersTable.GrantReadData(readOrderHandler)
	ordersTable.GrantReadWriteData(updateOrderHandler)
	ordersTable.GrantWriteData(deleteOrderHandler)

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
