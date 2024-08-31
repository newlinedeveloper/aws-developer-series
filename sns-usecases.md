### 1. **Fan-Out Architecture**

**Use Case:**
You have a system where an event (e.g., a new user registration) triggers multiple actions such as sending a welcome email, updating a database, and logging the event. You can use SNS to fan out the event to multiple subscribers (Lambda functions, SQS queues, etc.).

**Concepts Demonstrated:**
- SNS Topic Creation
- Multiple Subscribers (Lambda, SQS)

**CDK Example:**

```go
// Create an SNS Topic
topic := awssns.NewTopic(stack, jsii.String("UserRegistrationTopic"), &awssns.TopicProps{
    DisplayName: jsii.String("User Registration Topic"),
})

// Create Lambda function for sending welcome email
sendWelcomeEmail := awslambda.NewFunction(stack, jsii.String("SendWelcomeEmailFunction"), &awslambda.FunctionProps{
    Runtime: awslambda.Runtime_NODEJS_14_X(),
    Code:    awslambda.Code_FromAsset(jsii.String("lambda")),
    Handler: jsii.String("sendWelcomeEmail.handler"),
})

// Create Lambda function for logging the event
logEvent := awslambda.NewFunction(stack, jsii.String("LogEventFunction"), &awslambda.FunctionProps{
    Runtime: awslambda.Runtime_NODEJS_14_X(),
    Code:    awslambda.Code_FromAsset(jsii.String("lambda")),
    Handler: jsii.String("logEvent.handler"),
})

// Subscribe Lambda functions to SNS Topic
topic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(sendWelcomeEmail))
topic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(logEvent))
```

**Explanation:**
- **UserRegistrationTopic**: The SNS topic is created to publish messages when a new user registers.
- **SendWelcomeEmailFunction** and **LogEventFunction**: Lambda functions are created to handle different actions triggered by the SNS topic.
- **Subscriptions**: The Lambda functions are subscribed to the SNS topic, enabling them to receive messages when the topic is published.

### 2. **Email Notifications for Critical Alerts**

**Use Case:**
You need to send email notifications for critical system alerts or errors. SNS can be used to trigger email notifications via Amazon Simple Email Service (SES) or other email endpoints.

**Concepts Demonstrated:**
- SNS Topic Creation
- Email Subscription

**CDK Example:**

```go
// Create an SNS Topic for critical alerts
criticalAlertTopic := awssns.NewTopic(stack, jsii.String("CriticalAlertTopic"), &awssns.TopicProps{
    DisplayName: jsii.String("Critical Alert Topic"),
})

// Subscribe an email address to the SNS Topic
criticalAlertTopic.AddSubscription(awssnssubscriptions.NewEmailSubscription(jsii.String("alert@example.com")))
```

**Explanation:**
- **CriticalAlertTopic**: An SNS topic is created for critical alerts.
- **EmailSubscription**: An email address (`alert@example.com`) is subscribed to the SNS topic, so it receives notifications when messages are published.

### 3. **Push Notifications to Mobile Devices**

**Use Case:**
You want to send push notifications to mobile devices (iOS, Android) when certain events occur, such as a new message or an app update.

**Concepts Demonstrated:**
- SNS Platform Application for Mobile Push
- SNS Topic Creation
- Platform Endpoint Subscription

**CDK Example:**

```go
// Create an SNS Platform Application for iOS
platformApp := awssns.NewCfnPlatformApplication(stack, jsii.String("MyIosPlatformApp"), &awssns.CfnPlatformApplicationProps{
    Platform:           jsii.String("APNS"), // Use "APNS_SANDBOX" for development
    PlatformCredential: jsii.String("your-ios-push-cert"),
})

// Create an SNS Topic for mobile push notifications
pushNotificationTopic := awssns.NewTopic(stack, jsii.String("PushNotificationTopic"), &awssns.TopicProps{
    DisplayName: jsii.String("Push Notification Topic"),
})

// Create an endpoint for a mobile device
platformEndpoint := awssns.NewCfnPlatformEndpoint(stack, jsii.String("MyIosDeviceEndpoint"), &awssns.CfnPlatformEndpointProps{
    PlatformApplicationArn: platformApp.AttrArn(),
    Token:                  jsii.String("device-token"),
})

// Subscribe the mobile endpoint to the topic
pushNotificationTopic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(awssns.NewCfnSubscription(stack, jsii.String("MyIosDeviceSubscription"), &awssns.CfnSubscriptionProps{
    Endpoint: platformEndpoint.Ref(),
    Protocol: jsii.String("application"),
    TopicArn: pushNotificationTopic.TopicArn(),
})))
```

**Explanation:**
- **MyIosPlatformApp**: An SNS platform application is created for iOS push notifications.
- **PushNotificationTopic**: An SNS topic is created for publishing push notifications.
- **MyIosDeviceEndpoint**: An endpoint is created for a specific mobile device using its device token.
- **MyIosDeviceSubscription**: The mobile endpoint is subscribed to the SNS topic, enabling the device to receive push notifications.

### 4. **S3 Event Notifications with SNS**

**Use Case:**
You want to notify multiple systems or users when an object is created or deleted in an S3 bucket. SNS can be used to broadcast S3 events to multiple subscribers.

**Concepts Demonstrated:**
- S3 Event Notifications
- SNS Topic Creation
- SNS Subscription for S3 Events

**CDK Example:**

```go
// Create an SNS Topic for S3 event notifications
s3EventTopic := awssns.NewTopic(stack, jsii.String("S3EventTopic"), &awssns.TopicProps{
    DisplayName: jsii.String("S3 Event Topic"),
})

// Create an S3 bucket
bucket := awss3.NewBucket(stack, jsii.String("MyBucket"), &awss3.BucketProps{
    Versioned: jsii.Bool(true),
})

// Set up S3 event notification to SNS Topic
bucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3notifications.NewSnsDestination(s3EventTopic))

// Subscribe an email to the SNS Topic to receive S3 event notifications
s3EventTopic.AddSubscription(awssnssubscriptions.NewEmailSubscription(jsii.String("admin@example.com")))
```

**Explanation:**
- **S3EventTopic**: An SNS topic is created to publish notifications about S3 events.
- **MyBucket**: An S3 bucket is created and configured to trigger an SNS notification when an object is created.
- **EmailSubscription**: An email address is subscribed to receive notifications when objects are created in the S3 bucket.

### 5. **Cross-Account Notifications**

**Use Case:**
You need to send notifications from resources in one AWS account to another. SNS can be used to publish messages to topics in one account, and subscribers in another account can receive them.

**Concepts Demonstrated:**
- Cross-Account SNS Topic Subscription
- SNS Topic Policy

**CDK Example:**

```go
// Create an SNS Topic in Account A
crossAccountTopic := awssns.NewTopic(stack, jsii.String("CrossAccountTopic"), &awssns.TopicProps{
    DisplayName: jsii.String("Cross-Account Topic"),
})

// Define a policy to allow another AWS account to subscribe to this topic
crossAccountTopic.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
    Actions:   &[]*string{jsii.String("sns:Subscribe")},
    Principals: &[]awsiam.IPrincipal{awsiam.NewAccountPrincipal(jsii.String("another-account-id"))},
    Resources: &[]*string{crossAccountTopic.TopicArn()},
}))

// In Account B, subscribe to the topic (this part would be in a different CDK app in the other account)
awssns.NewCfnSubscription(stack, jsii.String("CrossAccountSubscription"), &awssns.CfnSubscriptionProps{
    Endpoint: jsii.String("target-endpoint@example.com"), // e.g., email or Lambda
    Protocol: jsii.String("email"), // Could be "lambda", "sqs", etc.
    TopicArn: jsii.String("arn:aws:sns:region:account-id:CrossAccountTopic"), // Replace with actual ARN
})
```

**Explanation:**
- **CrossAccountTopic**: An SNS topic is created in one account (Account A) and configured to allow subscriptions from another account (Account B).
- **Policy Statement**: A resource policy is added to the topic to grant `Subscribe` permission to another AWS account.
- **CrossAccountSubscription**: In Account B, a subscription is created to the SNS topic in Account A, allowing cross-account communication.

### 6. **Event-Driven Microservices Communication**

**Use Case:**
You are building a microservices architecture where services need to communicate asynchronously. SNS can be used to publish events that multiple microservices can subscribe to, allowing for loose coupling and scalability.

**Concepts Demonstrated:**
- SNS Topic Creation
- Multiple Microservice Subscriptions
- Event-Driven Architecture

**CDK Example:**

```go
// Create an SNS Topic for microservices communication
microservicesTopic := awssns.NewTopic(stack, jsii.String("MicroservicesTopic"), &awssns.TopicProps{
    DisplayName: jsii.String("Microservices Communication Topic"),
})

// Create Lambda functions for different microservices
serviceA := awslambda.NewFunction(stack, jsii.String("ServiceAFunction"), &awsl

ambda.FunctionProps{
    Runtime: awslambda.Runtime_GO_1_X(),
    Code:    awslambda.Code_FromAsset(jsii.String("lambda")),
    Handler: jsii.String("serviceA.handler"),
})

serviceB := awslambda.NewFunction(stack, jsii.String("ServiceBFunction"), &awslambda.FunctionProps{
    Runtime: awslambda.Runtime_GO_1_X(),
    Code:    awslambda.Code_FromAsset(jsii.String("lambda")),
    Handler: jsii.String("serviceB.handler"),
})

// Subscribe microservices to the SNS Topic
microservicesTopic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(serviceA))
microservicesTopic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(serviceB))
```

**Explanation:**
- **MicroservicesTopic**: An SNS topic is created to enable communication between microservices.
- **ServiceAFunction** and **ServiceBFunction**: Lambda functions represent different microservices that subscribe to the SNS topic to receive events.
- **Subscriptions**: Each microservice is subscribed to the SNS topic, enabling them to receive and react to events published by other services.

### Conclusion

These use cases demonstrate the versatility of SNS in various scenarios, from event-driven architectures to cross-account notifications. By leveraging SNS with AWS CDK in Golang, you can build scalable, decoupled, and resilient applications.