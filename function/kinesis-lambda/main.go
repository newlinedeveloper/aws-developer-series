package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, kinesisEvent events.KinesisEvent) (string, error) {
	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		data := string(kinesisRecord.Data)

		// Process each record (for demonstration, we're just printing the data)
		log.Printf("Processing record with data: %s", data)
	}

	return "Processed Kinesis records", nil
}

func main() {
	lambda.Start(handleRequest)
}
