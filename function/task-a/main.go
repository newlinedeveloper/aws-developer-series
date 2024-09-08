package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

// Response from Lambda A
type TaskAResponse struct {
	Message string `json:"message"`
	Data    struct {
		Value  int    `json:"value"`
		Status string `json:"status"`
	} `json:"data"`
}

// Lambda A handler function
func HandleRequest(ctx context.Context) (TaskAResponse, error) {
	log.Println("Lambda A executed")

	// Simulate processing
	taskAResult := TaskAResponse{
		Message: "Task A completed successfully",
	}
	taskAResult.Data.Value = 42
	taskAResult.Data.Status = "success"

	// Log the output
	result, _ := json.Marshal(taskAResult)
	log.Println("Task A Result: ", string(result))

	// Return the result
	return taskAResult, nil
}

func main() {
	lambda.Start(HandleRequest)
}
