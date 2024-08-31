package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
)

type Order struct {
	OrderID       string `json:"OrderID"`
	OrderDateTime string `json:"OrderDateTime"`
	Product       string `json:"Product"`
	Quantity      int    `json:"Quantity"`
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.GetItemInput{
		TableName: jsii.String("OrdersTable"),
		Key: map[string]types.AttributeValue{
			"OrderID":       &types.AttributeValueMemberS{Value: "1"},
			"OrderDateTime": &types.AttributeValueMemberS{Value: "2024-07-30T12:00:00Z"},
		},
	}

	result, err := svc.GetItem(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to read item: %v\n", err)
		return
	}

	// Unmarshal the result into the Order struct
	var order Order
	err = attributevalue.UnmarshalMap(result.Item, &order)
	if err != nil {
		fmt.Printf("Failed to unmarshal item: %v\n", err)
		return
	}

	// Print the unmarshaled order in a readable format
	prettyOrder, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal order: %v\n", err)
		return
	}

	fmt.Println("Order Details:")
	fmt.Println(string(prettyOrder))
}
