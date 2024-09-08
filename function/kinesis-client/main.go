package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func main() {
	// Create an AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))

	// Create Kinesis client
	svc := kinesis.New(sess)

	streamName := "my-data-stream"

	// Sample data to send
	data := []byte("Welcome to AWS Developer series")

	// Put the record into the Kinesis stream
	_, err := svc.PutRecord(&kinesis.PutRecordInput{
		Data:         data,
		PartitionKey: aws.String("partitionKey-1"), // Partition key
		StreamName:   aws.String(streamName),
	})

	if err != nil {
		fmt.Println("Error putting data to stream:", err)
		return
	}

	fmt.Println("Successfully sent data to stream!")
}
