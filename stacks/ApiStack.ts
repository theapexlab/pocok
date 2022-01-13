import {
  Stack,
  App,
  StackProps,
  Api,
  Queue,
} from "@serverless-stack/resources";
import { QueueStack } from "./QueueStack";
import { StorageStack } from "./StorageStack";

interface AdditionalStackProps {
  storageStack: StorageStack;
  queueStack: QueueStack;
}

export class ApiStack extends Stack {
  constructor(
    scope: App,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    const api = new Api(this, "Api", {
      routes: {
        "POST /webhooks/pipedream": {
          function: {
            handler: "src/api/process_email/main.go",
            environment: {
              queueUrl: additionalStackProps?.queueStack.uploadQueue.sqsQueue
                .queueUrl as string,
            },
            permissions: [additionalStackProps?.queueStack.uploadQueue as Queue],
          },
        },
      },
    });

    this.addOutputs({
      ApiEndpoint: api.url,
    });
  }
}
