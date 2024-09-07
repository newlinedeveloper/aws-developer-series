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

type SNSMessage struct {
	Message string `json:"Message"`
}

// Handler function to process SQS messages
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	fmt.Println("SQS event processing --------------")
	for _, sqsMessage := range sqsEvent.Records {
		// The SQS message body contains an SNS message
		fmt.Println("Message => ", sqsMessage)
		fmt.Printf("SQS Message Body: %s\n", sqsMessage.Body)

		// Unmarshal the SNS message from the SQS message body
		var snsMessage SNSMessage
		err := json.Unmarshal([]byte(sqsMessage.Body), &snsMessage)
		if err != nil {
			return fmt.Errorf("failed to unmarshal SQS message as SNS message: %v", err)
		}

		// The SNS message contains the actual S3 event
		fmt.Printf("SNS Message: %s\n", snsMessage.Message)

		// Unmarshal the S3 event from the SNS message
		var s3Event S3Event
		err = json.Unmarshal([]byte(snsMessage.Message), &s3Event)
		if err != nil {
			return fmt.Errorf("failed to unmarshal SNS message into S3 event: %v", err)
		}

		// Process each S3 event record
		for _, s3Record := range s3Event.Records {
			bucketName := s3Record.S3.Bucket.Name
			objectKey := s3Record.S3.Object.Key
			objectSize := s3Record.S3.Object.Size

			fmt.Printf("File uploaded: %s/%s (Size: %d bytes)\n", bucketName, objectKey, objectSize)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
