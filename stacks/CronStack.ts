import { Construct } from "@aws-cdk/core";
import { Cron, Stack, StackProps, Queue } from "@serverless-stack/resources";
import { QueueStack } from "./QueueStack";

type AdditionalStackProps = {
  queueStack: QueueStack;
};

export class CronStack extends Stack {
  emailCron: Cron;

  constructor(
    scope: Construct,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    this.emailCron = new Cron(this, "EmailCron", {
      schedule: "cron(0 15 * * ? *)",
      job: {
        function: {
          handler: "src/cron/invoice_summary/main.go",
          environment: {
            queueUrl: additionalStackProps?.queueStack.emailSenderQueue.sqsQueue
              .queueUrl as string,
          },
          permissions: [
            additionalStackProps?.queueStack.emailSenderQueue as Queue,
          ],
        },
      },
    });
  }
}
