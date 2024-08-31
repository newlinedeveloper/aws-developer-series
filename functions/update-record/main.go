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

	// Extract quantity from the request body
	var requestBody struct {
		Quantity string `json:"quantity"`
	}
	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, fmt.Errorf("invalid request body: %v", err)
	}

	// Construct the update input
	input := &dynamodb.UpdateItemInput{
		TableName: jsii.String("OrdersTable"),
		Key: map[string]types.AttributeValue{
			"OrderID": &types.AttributeValueMemberS{Value: orderID},
		},
		UpdateExpression:          jsii.String("SET Quantity = :q"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":q": &types.AttributeValueMemberN{Value: requestBody.Quantity}},
	}

	_, err = svc.UpdateItem(ctx, input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to update item",
		}, fmt.Errorf("failed to update item: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Successfully updated item",
	}, nil
}

func main() {
	lambda.Start(handler)
}
