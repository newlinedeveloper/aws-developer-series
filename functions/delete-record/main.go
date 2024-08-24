package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	orderID := request.PathParameters["order_id"]

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to load config: %v", err),
		}, nil
	}

	svc := dynamodb.NewFromConfig(cfg)

	// Retrieve OrderDateTime from DynamoDB (Assuming it's needed)
	getItemInput := &dynamodb.GetItemInput{
		TableName: jsii.String(os.Getenv("ORDERS_TABLE")),
		Key: map[string]types.AttributeValue{
			"OrderID": &types.AttributeValueMemberS{Value: orderID},
		},
	}

	getItemOutput, err := svc.GetItem(ctx, getItemInput)
	if err != nil || getItemOutput.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("Order with ID %s not found", orderID),
		}, nil
	}

	orderDateTime := getItemOutput.Item["OrderDateTime"].(*types.AttributeValueMemberS).Value

	// Delete the item from DynamoDB
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: jsii.String(os.Getenv("ORDERS_TABLE")),
		Key: map[string]types.AttributeValue{
			"OrderID":       &types.AttributeValueMemberS{Value: orderID},
			"OrderDateTime": &types.AttributeValueMemberS{Value: orderDateTime},
		},
	}

	_, err = svc.DeleteItem(ctx, deleteItemInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to delete item: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Successfully deleted order with ID %s", orderID),
	}, nil
}

func main() {
	lambda.Start(handler)
}
