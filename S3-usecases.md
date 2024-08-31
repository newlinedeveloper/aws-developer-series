### 1. **Static Website Hosting with S3**

**Use Case:**
You need to host a static website (HTML, CSS, JavaScript files) on AWS. S3 provides a cost-effective solution for this.

**Concepts Demonstrated:**
- S3 Bucket
- Public Access
- Static Website Hosting

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("WebsiteBucket"), &awss3.BucketProps{
    WebsiteIndexDocument: jsii.String("index.html"),
    PublicReadAccess:     jsii.Bool(true),
    RemovalPolicy:        awscdk.RemovalPolicy_DESTROY,
})

// Output the website URL
awscdk.NewCfnOutput(stack, jsii.String("WebsiteURL"), &awscdk.CfnOutputProps{
    Value: bucket.BucketWebsiteUrl(),
})
```

**Explanation:**
- The `WebsiteBucket` is configured to serve as a static website by setting `WebsiteIndexDocument`.
- `PublicReadAccess` is set to true to allow public access to the website content.
- `BucketWebsiteUrl` is outputted so you can easily find and access your hosted website.

### 2. **Secure Data Storage with S3**

**Use Case:**
You need to store sensitive data in S3 and ensure that only specific IAM roles or users can access it.

**Concepts Demonstrated:**
- S3 Bucket
- Bucket Policy
- Encryption

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("SecureBucket"), &awss3.BucketProps{
    Encryption:           awss3.BucketEncryption_S3_MANAGED,
    BlockPublicAccess:    awss3.BlockPublicAccess_BLOCK_ALL(),
    RemovalPolicy:        awscdk.RemovalPolicy_DESTROY,
})

bucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
    Actions:   jsii.Strings("s3:GetObject"),
    Resources: jsii.Strings(fmt.Sprintf("%s/*", *bucket.BucketArn())),
    Principals: &[]awsiam.IPrincipal{
        awsiam.NewAccountRootPrincipal(),
    },
}))
```

**Explanation:**
- The `SecureBucket` uses S3-managed encryption (`S3_MANAGED`) to encrypt all objects stored in the bucket.
- `BlockPublicAccess_BLOCK_ALL` ensures that no public access is allowed.
- A custom bucket policy allows only the root account to access objects within the bucket.

### 3. **Versioned Storage for File Backups**

**Use Case:**
You need to store backups of files where you might need to retrieve or restore previous versions of the files.

**Concepts Demonstrated:**
- S3 Bucket Versioning

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("VersionedBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
})
```

**Explanation:**
- The `VersionedBucket` enables versioning, allowing you to retain multiple versions of objects. This is useful for backups where you might need to roll back to a previous state.

### 4. **Lifecycle Management for Data Archiving**

**Use Case:**
You want to move infrequently accessed data to a cheaper storage class (like Glacier) after a certain period and eventually delete it.

**Concepts Demonstrated:**
- S3 Bucket Lifecycle Rules
- Glacier Storage Class

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("LifecycleBucket"), &awss3.BucketProps{
    RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
    LifecycleRules: &[]*awss3.LifecycleRule{
        {
            Id:              jsii.String("MoveToGlacier"),
            Prefix:          jsii.String("backup/"),
            Transitions: []*awss3.Transition{
                {
                    StorageClass: awss3.StorageClass_GLACIER,
                    TransitionAfter: awscdk.Duration_Days(jsii.Number(30)),
                },
            },
            Expiration: awscdk.Duration_Days(jsii.Number(365)),
        },
    },
})
```

**Explanation:**
- The `LifecycleBucket` has a lifecycle rule that transitions objects under the `backup/` prefix to Glacier after 30 days and deletes them after 365 days.
- This is ideal for archival storage, where cost-effectiveness is key.

### 5. **Event-Driven Processing with S3 and Lambda**

**Use Case:**
You need to process files as soon as they are uploaded to an S3 bucket (e.g., image processing, data parsing).

**Concepts Demonstrated:**
- S3 Event Notifications
- Lambda Function Trigger

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("ProcessingBucket"), &awss3.BucketProps{
    RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
})

lambdaFunction := awslambda.NewFunction(stack, jsii.String("S3ProcessingFunction"), &awslambda.FunctionProps{
    Runtime: awslambda.Runtime_NODEJS_14_X(),
    Code:    awslambda.Code_FromAsset(jsii.String("lambda")),
    Handler: jsii.String("index.handler"),
})

bucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3.NewLambdaDestination(lambdaFunction))

bucket.GrantRead(lambdaFunction)
```

**Explanation:**
- The `ProcessingBucket` triggers a Lambda function whenever a new object is uploaded.
- The Lambda function can then process the uploaded file (e.g., resize an image, analyze data).
- `GrantRead` ensures the Lambda function has the necessary permissions to read from the bucket.

### 6. **Cross-Account Data Sharing with S3**

**Use Case:**
You need to share data stored in an S3 bucket with another AWS account securely.

**Concepts Demonstrated:**
- Bucket Policy
- Cross-Account Access

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("SharedBucket"), &awss3.BucketProps{
    RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
})

bucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
    Actions:   jsii.Strings("s3:GetObject"),
    Resources: jsii.Strings(fmt.Sprintf("%s/*", *bucket.BucketArn())),
    Principals: &[]awsiam.IPrincipal{
        awsiam.NewAccountPrincipal(jsii.String("123456789012")),  // Replace with the target account ID
    },
}))
```

**Explanation:**
- The `SharedBucket` allows another AWS account (specified by its account ID) to access objects within the bucket securely.
- This setup is useful when you need to share data across AWS accounts while maintaining security controls.


### 7. **Same-Region Data Replication**

**Use Case:**
You want to replicate data within the same AWS region across multiple S3 buckets. This can be useful for data redundancy, backup, or sharing data between different environments (e.g., development and production).

**Concepts Demonstrated:**
- S3 Replication (Same-Region)
- IAM Roles for Replication

**CDK Example:**

```go
// Source bucket
sourceBucket := awss3.NewBucket(stack, jsii.String("SourceBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
})

// Destination bucket (must be in the same region)
destinationBucket := awss3.NewBucket(stack, jsii.String("DestinationBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
})

// Create IAM role for replication
replicationRole := awsiam.NewRole(stack, jsii.String("ReplicationRole"), &awsiam.RoleProps{
    AssumedBy: awsiam.NewServicePrincipal(jsii.String("s3.amazonaws.com")),
})

// Grant the role permission to replicate objects
sourceBucket.GrantRead(replicationRole)
destinationBucket.GrantWrite(replicationRole)

// Set up the replication configuration
sourceBucket.AddReplicationRule(&awss3.ReplicationRule{
    Destination: &awss3.ReplicationDestination{
        Bucket: destinationBucket,
        Account: jsii.String("your-aws-account-id"),
    },
    Status:       awss3.ReplicationRuleStatus_ENABLED,
    Role:         replicationRole.RoleArn(),
})
```

**Explanation:**
- **Source and Destination Buckets:** Both buckets are created in the same AWS region.
- **Replication Role:** An IAM role is created to handle the replication process, with permissions to read from the source bucket and write to the destination bucket.
- **Replication Rule:** The rule specifies that objects from the source bucket should be replicated to the destination bucket.

### 8. **Cross-Region Data Replication**

**Use Case:**
You want to replicate data across different AWS regions for disaster recovery, compliance, or reducing latency by keeping data close to users.

**Concepts Demonstrated:**
- S3 Replication (Cross-Region)
- Cross-Region IAM Roles

**CDK Example:**

```go
// Source bucket in the primary region
sourceBucket := awss3.NewBucket(stack, jsii.String("SourceBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
})

// Destination bucket in a different region
destinationBucket := awss3.NewBucket(stack, jsii.String("DestinationBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
    BucketName:     jsii.String("destination-bucket-in-another-region"),
    Region:         jsii.String("us-west-2"),  // Specify the region
})

// Create IAM role for cross-region replication
replicationRole := awsiam.NewRole(stack, jsii.String("CrossRegionReplicationRole"), &awsiam.RoleProps{
    AssumedBy: awsiam.NewServicePrincipal(jsii.String("s3.amazonaws.com")),
})

// Grant the role permission to replicate objects
sourceBucket.GrantRead(replicationRole)
destinationBucket.GrantWrite(replicationRole)

// Set up the cross-region replication configuration
sourceBucket.AddReplicationRule(&awss3.ReplicationRule{
    Destination: &awss3.ReplicationDestination{
        Bucket: destinationBucket,
        Account: jsii.String("your-aws-account-id"),
    },
    Status:       awss3.ReplicationRuleStatus_ENABLED,
    Role:         replicationRole.RoleArn(),
})
```

**Explanation:**
- **Source and Destination Buckets:** The source bucket is in the primary region, while the destination bucket is in a different region (`us-west-2` in this example).
- **Cross-Region Replication Role:** An IAM role is created specifically for cross-region replication, ensuring that it has the appropriate permissions in both regions.
- **Cross-Region Replication Rule:** This rule configures the source bucket to replicate its objects to the destination bucket in another region.


AWS S3 is an extremely versatile service with a wide range of use cases beyond basic storage and replication. Here are some additional advanced use cases that you can implement using AWS S3, especially in conjunction with other AWS services:

### 9. **Data Lake for Big Data Analytics**

**Use Case:**
You want to build a data lake that stores massive amounts of structured and unstructured data for big data analytics, machine learning, and business intelligence.

**Concepts Demonstrated:**
- S3 Data Lake
- Integration with AWS Glue, Athena, and Redshift Spectrum

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("DataLakeBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
    BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
})

// Grant access to AWS Glue for cataloging data
bucket.GrantReadWrite(awsiam.NewRole(stack, jsii.String("GlueServiceRole"), &awsiam.RoleProps{
    AssumedBy: awsiam.NewServicePrincipal(jsii.String("glue.amazonaws.com")),
}))

// Grant access to AWS Athena for querying data
bucket.GrantRead(awsiam.NewRole(stack, jsii.String("AthenaServiceRole"), &awsiam.RoleProps{
    AssumedBy: awsiam.NewServicePrincipal(jsii.String("athena.amazonaws.com")),
}))
```

**Explanation:**
- The `DataLakeBucket` is used to store raw data that can be analyzed and processed by big data tools.
- AWS Glue can be used to catalog and prepare data, while Athena and Redshift Spectrum can query the data directly from S3.

### 10. **Serverless Web Application Backend**

**Use Case:**
You want to store media files, user uploads, or static content for a serverless web application where the backend is built using AWS Lambda, API Gateway, and DynamoDB.

**Concepts Demonstrated:**
- S3 Storage for Media Files
- Integration with Lambda and API Gateway

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("AppMediaBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
    BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
})

// Lambda function to handle file uploads
lambdaFunction := awslambda.NewFunction(stack, jsii.String("MediaUploadHandler"), &awslambda.FunctionProps{
    Runtime: awslambda.Runtime_NODEJS_14_X(),
    Code:    awslambda.Code_FromAsset(jsii.String("lambda")),
    Handler: jsii.String("index.handler"),
})

// API Gateway to expose Lambda function as a REST API
api := awsapigateway.NewLambdaRestApi(stack, jsii.String("AppApi"), &awsapigateway.LambdaRestApiProps{
    Handler: lambdaFunction,
})

// Grant Lambda function permission to write to the S3 bucket
bucket.GrantWrite(lambdaFunction)
```

**Explanation:**
- The `AppMediaBucket` is used to store user-generated content or static files.
- A Lambda function handles file uploads, and API Gateway exposes this function as an API endpoint.

### 11. **Automated Backups and Data Archiving**

**Use Case:**
You want to automate backups of critical data and archive it to S3 for long-term storage, using S3 Glacier for cost-effective archiving.

**Concepts Demonstrated:**
- S3 Lifecycle Rules
- S3 Glacier for Archiving

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("BackupBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
    LifecycleRules: &[]*awss3.LifecycleRule{
        {
            Id:       jsii.String("ArchiveOldVersions"),
            Transitions: []*awss3.Transition{
                {
                    StorageClass: awss3.StorageClass_GLACIER,
                    TransitionAfter: awscdk.Duration_Days(jsii.Number(30)),
                },
            },
            NoncurrentVersionTransitions: []*awss3.NoncurrentVersionTransition{
                {
                    StorageClass: awss3.StorageClass_GLACIER,
                    TransitionAfter: awscdk.Duration_Days(jsii.Number(30)),
                },
            },
        },
    },
})
```

**Explanation:**
- The `BackupBucket` is configured with lifecycle rules to automatically transition older objects and noncurrent versions to S3 Glacier for long-term storage and cost savings.

### 12. **Content Distribution with CloudFront**

**Use Case:**
You want to distribute content (e.g., images, videos, or web assets) globally with low latency using Amazon CloudFront, with S3 as the origin.

**Concepts Demonstrated:**
- S3 as an Origin for CloudFront
- Integration with CloudFront

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("ContentBucket"), &awss3.BucketProps{
    Versioned:      jsii.Bool(true),
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
})

cdn := awscloudfront.NewDistribution(stack, jsii.String("MyDistribution"), &awscloudfront.DistributionProps{
    DefaultBehavior: &awscloudfront.BehaviorOptions{
        Origin: awscloudfront.NewS3Origin(bucket),
    },
})

awscdk.NewCfnOutput(stack, jsii.String("DistributionDomainName"), &awscdk.CfnOutputProps{
    Value: cdn.DistributionDomainName(),
})
```

**Explanation:**
- The `ContentBucket` serves as the origin for a CloudFront distribution.
- CloudFront caches content globally, reducing latency for users around the world and improving the performance of your web applications.

### 13. **Log Aggregation and Analysis**

**Use Case:**
You want to aggregate logs from different services (e.g., web servers, applications) into an S3 bucket for centralized storage and analysis using Amazon Athena.

**Concepts Demonstrated:**
- Centralized Log Storage
- Integration with Athena for Analysis

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("LogBucket"), &awss3.BucketProps{
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
    LifecycleRules: &[]*awss3.LifecycleRule{
        {
            Id:             jsii.String("ArchiveLogs"),
            Transitions: []*awss3.Transition{
                {
                    StorageClass: awss3.StorageClass_GLACIER,
                    TransitionAfter: awscdk.Duration_Days(jsii.Number(90)),
                },
            },
            Expiration: awscdk.Duration_Days(jsii.Number(365)),
        },
    },
})

// Grant Athena access to the logs for querying
bucket.GrantRead(awsiam.NewRole(stack, jsii.String("AthenaLogAccessRole"), &awsiam.RoleProps{
    AssumedBy: awsiam.NewServicePrincipal(jsii.String("athena.amazonaws.com")),
}))
```

**Explanation:**
- The `LogBucket` aggregates logs from various sources, with lifecycle rules to archive older logs to S3 Glacier.
- AWS Athena can query these logs directly from S3, allowing for powerful and cost-effective log analysis.

### 14. **Media Transcoding Pipeline**

**Use Case:**
You want to create a media processing pipeline where video files uploaded to S3 are automatically transcoded into different formats and resolutions using AWS Elastic Transcoder or AWS Elemental MediaConvert.

**Concepts Demonstrated:**
- S3 Event Notifications
- Integration with AWS Elastic Transcoder or MediaConvert

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("MediaBucket"), &awss3.BucketProps{
    RemovalPolicy:  awscdk.RemovalPolicy_DESTROY,
})

lambdaFunction := awslambda.NewFunction(stack, jsii.String("TranscoderTrigger"), &awslambda.FunctionProps{
    Runtime: awslambda.Runtime_NODEJS_14_X(),
    Code:    awslambda.Code_FromAsset(jsii.String("lambda")),
    Handler: jsii.String("index.handler"),
})

// Trigger Lambda function on new uploads to start transcoding
bucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3.NewLambdaDestination(lambdaFunction))

bucket.GrantReadWrite(lambdaFunction)
```

**Explanation:**
- The `MediaBucket` stores raw video files, which trigger a Lambda function on upload.
- The Lambda function could then start a media transcoding job using Elastic Transcoder or MediaConvert to process the video into different formats.


### 15. **Static Website Hosting**

**Use Case:**
You want to host a static website (e.g., HTML, CSS, JavaScript, images) directly from an S3 bucket. This is a cost-effective and scalable way to serve static content, especially for personal websites, documentation, or landing pages.

**Concepts Demonstrated:**
- S3 Static Website Hosting
- Configuring Bucket Policy for Public Access
- Custom Domain with Route 53 (Optional)

**CDK Example:**

```go
bucket := awss3.NewBucket(stack, jsii.String("WebsiteBucket"), &awss3.BucketProps{
    WebsiteIndexDocument: jsii.String("index.html"),
    WebsiteErrorDocument: jsii.String("error.html"),
    PublicReadAccess:     jsii.Bool(true),
    RemovalPolicy:        awscdk.RemovalPolicy_DESTROY,
})

// Output the website URL
awscdk.NewCfnOutput(stack, jsii.String("WebsiteURL"), &awscdk.CfnOutputProps{
    Value: bucket.BucketWebsiteUrl(),
})
```

**Explanation:**
- **WebsiteBucket**: The S3 bucket is configured for static website hosting, with `index.html` as the default document and `error.html` for error handling.
- **Public Read Access**: The bucket is made publicly accessible so that anyone can view the website.
- **Website URL Output**: The URL of the hosted website is outputted for easy access.

**Optional Enhancement:**

If you want to use a custom domain for your static website, you can integrate the S3 bucket with Route 53 for DNS management and optionally use CloudFront for SSL/TLS encryption.

```go
// Route 53 hosted zone
zone := awsroute53.HostedZone_FromLookup(stack, jsii.String("MyHostedZone"), &awsroute53.HostedZoneProviderProps{
    DomainName: jsii.String("example.com"),
})

// CloudFront distribution for custom domain and SSL/TLS
distribution := awscloudfront.NewDistribution(stack, jsii.String("WebsiteDistribution"), &awscloudfront.DistributionProps{
    DefaultBehavior: &awscloudfront.BehaviorOptions{
        Origin: awscloudfront.NewS3Origin(bucket),
    },
    DomainNames: &[]*string{
        jsii.String("www.example.com"),
    },
    Certificate: awscertificatemanager.NewCertificate(stack, jsii.String("WebsiteCertificate"), &awscertificatemanager.CertificateProps{
        DomainName:   jsii.String("www.example.com"),
        Validation:   awscertificatemanager.CertificateValidation_FromDns(zone),
    }),
})

// Route 53 record set for the custom domain
awsroute53.NewARecord(stack, jsii.String("WebsiteAliasRecord"), &awsroute53.ARecordProps{
    Zone:       zone,
    Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewCloudFrontTarget(distribution)),
    RecordName: jsii.String("www"),
})
```

**Explanation:**
- **Route 53 Hosted Zone**: Look up the hosted zone for your custom domain (e.g., `example.com`).
- **CloudFront Distribution**: Use CloudFront to distribute your static website globally, with SSL/TLS encryption using an ACM certificate.
- **Route 53 Record Set**: Create an alias record that points your custom domain (`www.example.com`) to the CloudFront distribution.
