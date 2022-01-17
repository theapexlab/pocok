import { Construct } from "@aws-cdk/core";
import { Cron, Stack, StackProps, Table } from "@serverless-stack/resources";
import { QueueStack } from "./QueueStack";
import { StorageStack } from "./StorageStack";

type AdditionalStackProps = {
  storageStack: StorageStack;
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
      schedule: "cron(0 16 * * ? *)",
      job: {
        function: {
          handler: "src/cron/invoice_summary/main.go",
          environment: {
            tableName: additionalStackProps?.storageStack.invoiceTable
              .tableName as string,
            queueUrl: additionalStackProps?.queueStack.emailSenderQueue.sqsQueue
              .queueUrl as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceTable as Table,
          ],
        },
      },
    });
  }
}
