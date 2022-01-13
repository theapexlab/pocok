import { App } from "@serverless-stack/resources";
import { ApiStack } from "./ApiStack";
import { StorageStack } from "./StorageStack";
import { QueueStack } from "./QueueStack";

export default function main(app: App): void {
  // Set default runtime for all functions
  app.setDefaultFunctionProps({
    runtime: "go1.x",
  });

  const storageStack = new StorageStack(app, "storage-stack");
  const queueStack = new QueueStack(app, "queue-stack", {}, { storageStack });

  new ApiStack(app, "api-stack", {}, { storageStack, queueStack });
}
