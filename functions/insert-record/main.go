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

type Request struct {
	OrderID       string `json:"order_id"`
	OrderDateTime string `json:"order_datetime"`
	Product       string `json:"product"`
	Quantity      string `json:"quantity"`
}

func HandleRequest(ctx context.Context, req Request) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	item := map[string]types.AttributeValue{
		"OrderID":       &types.AttributeValueMemberS{Value: req.OrderID},
		"OrderDateTime": &types.AttributeValueMemberS{Value: req.OrderDateTime},
		"Product":       &types.AttributeValueMemberS{Value: req.Product},
		"Quantity":      &types.AttributeValueMemberN{Value: req.Quantity},
	}

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: jsii.String("OrdersTable"),
		Item:      item,
	})

	if err != nil {
		return "", fmt.Errorf("failed to insert item: %v", err)
	}

	return "Successfully inserted item", nil
}

func main() {
	lambda.Start(HandleRequest)
}
