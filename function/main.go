package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// S3Detail represents the S3 object information sent via EventBridge
type S3Detail struct {
	Bucket struct {
		Name string `json:"name"`
	} `json:"bucket"`
	Object struct {
		Key  string `json:"key"`
		Size int64  `json:"size"`
	} `json:"object"`
}

// Handler function to process EventBridge messages containing S3 events
func handler(ctx context.Context, event events.CloudWatchEvent) error {
	fmt.Println("EventBridge event received --------------")

	// Unmarshal the S3 event details from the EventBridge event
	var s3Detail S3Detail
	err := json.Unmarshal(event.Detail, &s3Detail)
	if err != nil {
		return fmt.Errorf("failed to unmarshal EventBridge message into S3 detail: %v", err)
	}

	// Extract bucket name, object key, and size
	bucketName := s3Detail.Bucket.Name
	objectKey := s3Detail.Object.Key
	objectSize := s3Detail.Object.Size

	// Log the details of the uploaded file
	fmt.Printf("File uploaded: %s/%s (Size: %d bytes)\n", bucketName, objectKey, objectSize)

	return nil
}

func main() {
	lambda.Start(handler)
}
