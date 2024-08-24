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
}

func handler(ctx context.Context, event MyEvent) (map[string]types.AttributeValue, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.GetItemInput{
		TableName: jsii.String("OrdersTable"),
		Key: map[string]types.AttributeValue{
			"OrderID":       &types.AttributeValueMemberS{Value: event.OrderID},
			"OrderDateTime": &types.AttributeValueMemberS{Value: event.OrderDateTime},
		},
	}

	result, err := svc.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to read item: %w", err)
	}

	return result.Item, nil
}

func main() {
	lambda.Start(handler)
}
