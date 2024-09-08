package main

import (
	"developer-series/config"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctionstasks"
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
	taskAFunc := awslambda.NewFunction(stack, jsii.String(config.TaskAFunctionName), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-" + config.TaskAFunctionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(config.MemorySize),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(config.MaxDuration)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String(config.TaskACodePath), nil),
		Handler:      jsii.String(config.Handler),
	})

	taskBFunc := awslambda.NewFunction(stack, jsii.String(config.TaskBFunctionName), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-" + config.TaskBFunctionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(config.MemorySize),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(config.MaxDuration)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String(config.TaskBCodePath), nil),
		Handler:      jsii.String(config.Handler),
	})

	// Step Function task for Lambda A
	taskA := awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("InvokeLambdaA"), &awsstepfunctionstasks.LambdaInvokeProps{
		LambdaFunction: taskAFunc,
		OutputPath:     jsii.String("$.Payload"),
	})

	// Step Function task for Lambda B
	taskB := awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("InvokeLambdaB"), &awsstepfunctionstasks.LambdaInvokeProps{
		LambdaFunction: taskBFunc,
		OutputPath:     jsii.String("$.Payload"),
	})

	// Define a success state that Lambda A leads to Lambda B
	workflowChain := taskA.Next(taskB)

	// Create the state machine
	stateMachine := awsstepfunctions.NewStateMachine(stack, jsii.String("MyStateMachine"), &awsstepfunctions.StateMachineProps{
		Definition: workflowChain,
		Timeout:    awscdk.Duration_Minutes(jsii.Number(5)),
	})

	// Output the state machine ARN
	awscdk.NewCfnOutput(stack, jsii.String("StateMachineARN"), &awscdk.CfnOutputProps{
		Value: stateMachine.StateMachineArn(),
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
