package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
)

// Request defines the expected input structure from API Gateway
type Request struct {
	OrderID  string `json:"order_id"`
	Product  string `json:"product"`
	Quantity string `json:"quantity"` // Assuming quantity is sent as a string
}

// Response defines the structure of the response returned to API Gateway
type Response struct {
	Message string `json:"message"`
}

// HandleRequest handles the incoming request from API Gateway and processes it
func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to load SDK config",
		}, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	// Parse the request body
	var requestBody Request
	err = json.Unmarshal([]byte(req.Body), &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, fmt.Errorf("failed to unmarshal request body: %v", err)
	}

	// Generate the current timestamp for OrderDateTime
	orderDateTime := time.Now().UTC().Format(time.RFC3339)

	// Convert Quantity to int
	quantity, err := strconv.Atoi(requestBody.Quantity)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid quantity format",
		}, fmt.Errorf("failed to convert quantity to int: %v", err)
	}

	item := map[string]types.AttributeValue{
		"OrderID":       &types.AttributeValueMemberS{Value: requestBody.OrderID},
		"OrderDateTime": &types.AttributeValueMemberS{Value: orderDateTime},
		"Product":       &types.AttributeValueMemberS{Value: requestBody.Product},
		"Quantity":      &types.AttributeValueMemberN{Value: strconv.Itoa(quantity)},
	}

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: jsii.String("OrdersTable"),
		Item:      item,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to insert item",
		}, fmt.Errorf("failed to insert item: %v", err)
	}

	// Return a successful response
	responseBody := map[string]string{"message": "Successfully inserted item"}
	responseBodyJSON, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal response",
		}, fmt.Errorf("failed to marshal response: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBodyJSON),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
