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
  --username user@example.com \
  --password 'Password123!'


```

**Authenticate: Obtain a JWT token:**

```
aws cognito-idp initiate-auth \
  --client-id <user-pool-client-id> \
  --auth-flow USER_PASSWORD_AUTH \
  --auth-parameters USERNAME=user@example.com,PASSWORD='Password123!'


```

**Use the id_token to authenticate API requests with curl or Postman:**
```
curl -X GET https://<api-id>.execute-api.<region>.amazonaws.com/prod/items \
-H "Authorization: Bearer <id_token>"

```