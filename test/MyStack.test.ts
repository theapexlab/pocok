import { expect, haveResource } from "@aws-cdk/assert";
import { App } from "@serverless-stack/resources";
import { ApiStack } from "../stacks/ApiStack";

test("Test Stack", () => {
  const app = new App();
  // WHEN
  const stack = new ApiStack(app, "test-stack");
  // THEN
  expect(stack).to(haveResource("AWS::Lambda::Function"));
});
