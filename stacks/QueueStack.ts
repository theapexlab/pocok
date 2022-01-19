import { Construct } from "@aws-cdk/core";
import {
  Bucket,
  Queue,
  Stack,
  StackProps,
  Table,
} from "@serverless-stack/resources";
import { StorageStack } from "./StorageStack";

type AdditionalStackProps = {
  storageStack: StorageStack;
};

export class QueueStack extends Stack {
  queue: Queue;
  emailSenderQueue: Queue;

  constructor(
    scope: Construct,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    this.queue = new Queue(this, "Queue", {
      consumer: {
        function: {
          handler: "src/consumers/s3_uploader/main.go",
          environment: {
            tableName: additionalStackProps?.storageStack.invoiceTable
              .tableName as string,
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceTable as Table,
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
          ],
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });

    this.emailSenderQueue = new Queue(this, "EmailSender", {
      consumer: {
        function: {
          handler: "src/consumers/email_sender/main.go",
          environment: {
            domain: process.env.DOMAIN as string,
            sender: process.env.SENDER as string,
            apiKey: process.env.MAILGUN_API_KEY as string,
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            tableName: additionalStackProps?.storageStack.invoiceTable
              .tableName as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
            additionalStackProps?.storageStack.invoiceTable as Table,
          ],
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });
  }
}
