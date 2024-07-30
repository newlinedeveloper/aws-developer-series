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

	items := []map[string]types.AttributeValue{
		{
			"OrderID":       &types.AttributeValueMemberS{Value: "1"},
			"OrderDateTime": &types.AttributeValueMemberS{Value: "2024-07-30T12:00:00Z"},
			"Product":       &types.AttributeValueMemberS{Value: "Laptop"},
			"Quantity":      &types.AttributeValueMemberN{Value: "2"},
		},
		{
			"OrderID":       &types.AttributeValueMemberS{Value: "2"},
			"OrderDateTime": &types.AttributeValueMemberS{Value: "2024-07-30T13:00:00Z"},
			"Product":       &types.AttributeValueMemberS{Value: "Phone"},
			"Quantity":      &types.AttributeValueMemberN{Value: "5"},
		},
	}

	for _, item := range items {
		_, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: jsii.String("OrdersTable"),
			Item:      item,
		})
		if err != nil {
			fmt.Printf("Failed to insert item: %v\n", err)
		} else {
			fmt.Println("Successfully inserted item")
		}
	}
}
