package main

import (
	"fmt"
	"strconv"

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

	// Create a batch of 100 records to send
	var records []*kinesis.PutRecordsRequestEntry

	for i := 1; i <= 200; i++ {
		// Create a record
		data := []byte("Message " + strconv.Itoa(i) + ": Welcome to AWS Developer series")

		record := &kinesis.PutRecordsRequestEntry{
			Data:         data,
			PartitionKey: aws.String("partitionKey-" + strconv.Itoa(i)),
		}

		// Add the record to the batch
		records = append(records, record)
	}

	// Send the batch to Kinesis using PutRecords
	result, err := svc.PutRecords(&kinesis.PutRecordsInput{
		Records:    records,
		StreamName: aws.String(streamName),
	})

	if err != nil {
		fmt.Println("Error putting batch data to stream:", err)
		return
	}

	// Check if there are any failed records
	if *result.FailedRecordCount > 0 {
		fmt.Printf("Failed to put %d records to the stream\n", *result.FailedRecordCount)
	} else {
		fmt.Println("Successfully sent all records to stream!")
	}
}
