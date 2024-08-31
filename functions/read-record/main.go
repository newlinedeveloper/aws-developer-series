package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
)

// Handler function for API Gateway Proxy requests
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to load SDK config",
		}, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	// Extract order_id from the URL path parameter
	orderID, ok := req.PathParameters["order_id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing order_id path parameter",
		}, fmt.Errorf("missing order_id path parameter")
	}

	input := &dynamodb.GetItemInput{
		TableName: jsii.String("OrdersTable"),
		Key: map[string]types.AttributeValue{
			"OrderID": &types.AttributeValueMemberS{Value: orderID},
		},
	}

	result, err := svc.GetItem(ctx, input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to read item",
		}, fmt.Errorf("failed to read item: %v", err)
	}

	// Convert result.Item to a map[string]string
	item := make(map[string]string)
	for k, v := range result.Item {
		switch attr := v.(type) {
		case *types.AttributeValueMemberS:
			item[k] = attr.Value
		case *types.AttributeValueMemberN:
			item[k] = attr.Value
		default:
			// Handle other attribute types if needed
		}
	}

	// Convert item map to JSON
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal item",
		}, fmt.Errorf("failed to marshal item: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(itemJSON),
	}, nil
}

func main() {
	lambda.Start(handler)
}
