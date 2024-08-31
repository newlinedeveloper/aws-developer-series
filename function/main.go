package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3
		log.Printf("File uploaded: %s (size: %d bytes)", s3.Object.Key, s3.Object.Size)
	}

	fmt.Println("S3 event processed successfully")
}

func main() {
	lambda.Start(handler)
}
