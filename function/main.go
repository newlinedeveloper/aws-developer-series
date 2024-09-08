package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key  string `json:"key"`
				Size int64  `json:"size"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

// Handler function to process SQS messages containing S3 events
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	fmt.Println("SQS event processing --------------")
	for _, sqsMessage := range sqsEvent.Records {
		// Print the raw message for debugging
		fmt.Printf("SQS Message Body: %s\n", sqsMessage.Body)

		// Unmarshal the S3 event from the SQS message body
		var s3Event S3Event
		err := json.Unmarshal([]byte(sqsMessage.Body), &s3Event)
		if err != nil {
			return fmt.Errorf("failed to unmarshal SQS message into S3 event: %v", err)
		}

		// Process each S3 event record
		for _, s3Record := range s3Event.Records {
			bucketName := s3Record.S3.Bucket.Name
			objectKey := s3Record.S3.Object.Key
			objectSize := s3Record.S3.Object.Size

			// Log the details of the uploaded file
			fmt.Printf("File uploaded: %s/%s (Size: %d bytes)\n", bucketName, objectKey, objectSize)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
