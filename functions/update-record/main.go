package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
)

type MyEvent struct {
	OrderID       string `json:"order_id"`
	OrderDateTime string `json:"order_date_time"`
	Quantity      string `json:"quantity"`
}

func handler(ctx context.Context, event MyEvent) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.UpdateItemInput{
		TableName: jsii.String("OrdersTable"),
		Key: map[string]types.AttributeValue{
			"OrderID":       &types.AttributeValueMemberS{Value: event.OrderID},
			"OrderDateTime": &types.AttributeValueMemberS{Value: event.OrderDateTime},
		},
		UpdateExpression:          jsii.String("SET Quantity = :q"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":q": &types.AttributeValueMemberN{Value: event.Quantity}},
	}

	_, err = svc.UpdateItem(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to update item: %v", err)
	}

	return "Successfully updated item", nil
}

func main() {
	lambda.Start(handler)
}
