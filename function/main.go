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

func handler(ctx context.Context, snsEvent events.SNSEvent) error {
	fmt.Println("SNS notification event processing--------------")
	for _, record := range snsEvent.Records {
		snsMessage := record.SNS.Message
		var s3Event S3Event
		err := json.Unmarshal([]byte(snsMessage), &s3Event)
		if err != nil {
			return fmt.Errorf("failed to unmarshal SNS message: %v", err)
		}

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
