package main

import (
	"developer-series/config"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

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
