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

	input := &dynamodb.UpdateItemInput{
		TableName: jsii.String("OrdersTable"),
		Key: map[string]types.AttributeValue{
			"OrderID":       &types.AttributeValueMemberS{Value: "1"},
			"OrderDateTime": &types.AttributeValueMemberS{Value: "2024-07-30T12:00:00Z"},
		},
		UpdateExpression:          jsii.String("SET Quantity = :q"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":q": &types.AttributeValueMemberN{Value: "3"}},
	}

	_, err = svc.UpdateItem(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to update item: %v\n", err)
	} else {
		fmt.Println("Successfully updated item")
	}
}
