package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.DeleteItemInput{
		TableName: jsii.String("OrdersTable"),
		Key: map[string]types.AttributeValue{
			"OrderID":       &types.AttributeValueMemberS{Value: "1"},
			"OrderDateTime": &types.AttributeValueMemberS{Value: "2024-07-30T12:00:00Z"},
		},
	}

	_, err = svc.DeleteItem(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to delete item: %v\n", err)
	} else {
		fmt.Println("Successfully deleted item")
	}
}
