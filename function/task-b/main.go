package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

// Input for Lambda B (output from Lambda A)
type TaskBInput struct {
	Message string `json:"message"`
	Data    struct {
		Value  int    `json:"value"`
		Status string `json:"status"`
	} `json:"data"`
}

// Response from Lambda B
type TaskBResponse struct {
	Message        string      `json:"message"`
	ReceivedData   TaskBInput  `json:"receivedData"`
	AdditionalData interface{} `json:"additionalData"`
}

// Lambda B handler function
func HandleRequest(ctx context.Context, input TaskBInput) (TaskBResponse, error) {
	log.Println("Lambda B executed")
	log.Printf("Received input from Lambda A: %+v\n", input)

	// Simulate processing with received input
	taskBResult := TaskBResponse{
		Message:      "Task B completed successfully",
		ReceivedData: input,
		AdditionalData: map[string]string{
			"info":   "Task B processed the input",
			"status": "success",
		},
	}

	// Log the output
	result, _ := json.Marshal(taskBResult)
	log.Println("Task B Result: ", string(result))

	// Return the result
	return taskBResult, nil
}

func main() {
	lambda.Start(HandleRequest)
}
