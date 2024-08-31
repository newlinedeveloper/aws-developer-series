# Welcome to your CDK Go project!

This is a blank project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests


 ### Test the secure API

 **Sign Up: Create a user in the Cognito User Pool:**

```
aws cognito-idp sign-up \
  --client-id <user-pool-client-id> \
  --username testuser@gmail.com \
  --password 'Password123!'


```

**Admin confirms the user verification**

```
aws cognito-idp admin-confirm-sign-up \
  --user-pool-id <user-pool-id> \
  --username testuser@gmail.com


```

**Authenticate: Obtain a JWT token:**

```
aws cognito-idp initiate-auth \
  --client-id <user-pool-client-id> \
  --auth-flow USER_PASSWORD_AUTH \
  --auth-parameters USERNAME=testuser@gmail.com,PASSWORD='Password123!'


```

**Use the id_token to authenticate API requests with curl or Postman:**
```
curl -X GET https://<api-id>.execute-api.<region>.amazonaws.com/prod/items \
-H "Authorization: Bearer <id_token>"
```

**Create Order record payload**

```
{
  "order_id": "1",
  "product": "Laptop",
  "quantity": "2"
}

```