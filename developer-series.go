package main

import (
	"developer-series/config"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
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

	// Define the DynamoDB table
	ordersTable := awsdynamodb.NewTable(stack, jsii.String(config.TableName), &awsdynamodb.TableProps{
		TableName: jsii.String(config.TableName),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String(config.PartitionKey),
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

	// Create a Cognito User Pool
	userPool := awscognito.NewUserPool(stack, jsii.String("MyUserPool"), &awscognito.UserPoolProps{
		SelfSignUpEnabled: jsii.Bool(true),
		SignInAliases:     &awscognito.SignInAliases{Email: jsii.Bool(true)},
		AutoVerify:        &awscognito.AutoVerifiedAttrs{Email: jsii.Bool(true)},
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength:        jsii.Number(8),
			RequireLowercase: jsii.Bool(true),
			RequireUppercase: jsii.Bool(true),
			RequireDigits:    jsii.Bool(true),
			RequireSymbols:   jsii.Bool(true),
		},
		AccountRecovery: awscognito.AccountRecovery_EMAIL_ONLY,
		RemovalPolicy:   awscdk.RemovalPolicy_DESTROY,
	})

	_ = awscognito.NewUserPoolClient(stack, jsii.String("UserPoolClient"), &awscognito.UserPoolClientProps{
		UserPool:       userPool,
		GenerateSecret: jsii.Bool(false), // Set to true if you want to generate a client secret
		AuthFlows: &awscognito.AuthFlow{
			UserPassword: jsii.Bool(true), // Enable USER_PASSWORD_AUTH flow
		},
	})

	// Create API Gateway rest api.
	restApi := awsapigateway.NewRestApi(stack, jsii.String("LambdaRestApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(*stack.StackName() + "-LambdaRestApi"),
		Description: jsii.String("AWS Developer Series REST API"),
	})

	// Cognito Authorizer
	authorizer := awsapigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String("Authorizer"), &awsapigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{userPool},
	})

	// Define the base resource
	ordersApi := restApi.Root().AddResource(jsii.String("orders"), nil)

	// Define the "create" resource
	createOrderApiRes := ordersApi.AddResource(jsii.String("create"), nil)
	createOrderApiRes.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(createOrderHandler, nil), &awsapigateway.MethodOptions{
		Authorizer: authorizer,
	})

	// Define the "{order_id}" resource under "orders"
	orderIdResource := ordersApi.AddResource(jsii.String("{order_id}"), nil)

	// Define the GET method for the resource with path parameter
	orderIdResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(readOrderHandler, nil), &awsapigateway.MethodOptions{
		Authorizer: authorizer,
	})

	// Define the PUT method for the resource with path parameter
	orderIdResource.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(updateOrderHandler, nil), &awsapigateway.MethodOptions{
		Authorizer: authorizer,
	})

	// Define the DELETE method for the resource with path parameter
	orderIdResource.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(deleteOrderHandler, nil), &awsapigateway.MethodOptions{
		Authorizer: authorizer,
	})

	// Grant the Lambda function permissions to perform CRUD operations on the DynamoDB table
	ordersTable.GrantWriteData(createOrderHandler)
	ordersTable.GrantReadData(readOrderHandler)
	ordersTable.GrantReadWriteData(updateOrderHandler)
	ordersTable.GrantReadWriteData(deleteOrderHandler)

	return stack
}

func CreateLambdaHandler(stack awscdk.Stack, functionName string, codePath string) awslambda.Function {
	orderHandler := awslambda.NewFunction(stack, jsii.String(functionName), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-" + functionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(config.MemorySize),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(config.MaxDuration)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String(codePath), nil),
		Handler:      jsii.String(config.Handler),
	})

	return orderHandler
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
