import { expect, haveResource } from "@aws-cdk/assert";
import { App } from "@serverless-stack/resources";
import { ApiStack } from "../stacks/ApiStack";
import { QueueStack } from "../stacks/QueueStack";
import { StorageStack } from "../stacks/StorageStack";

test("Test Stack", () => {
  const app = new App();

  const queueStack = new QueueStack(app, "queue-stack");
  expect(queueStack).to(haveResource("AWS::SQS::Queue"));

  const storageStack = new StorageStack(app, "storage-stack");
  expect(storageStack).to(haveResource("AWS::S3::Bucket"));
  expect(storageStack).to(haveResource("AWS::DynamoDB::Table"));

  const apiStack = new ApiStack(
    app,
    "api-stack",
    {},
    {
      queueStack,
      storageStack,
    }
  );
  expect(apiStack).to(haveResource("AWS::ApiGateway::RestApi"));
  expect(apiStack).to(haveResource("AWS::ApiGateway::Resource"));
  expect(apiStack).to(haveResource("AWS::ApiGateway::Method"));
});
