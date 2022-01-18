import { Construct } from "@aws-cdk/core";
import { Bucket, Queue, Stack, StackProps } from "@serverless-stack/resources";
import { StorageStack } from "./StorageStack";

type AdditionalStackProps = {
  storageStack: StorageStack;
};

export class QueueStack extends Stack {
  invoiceQueue: Queue;
  processInvoiceQueue: Queue;

  constructor(
    scope: Construct,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    this.processInvoiceQueue =
      this.createProcessInvoiceQueue(additionalStackProps);
    this.invoiceQueue = this.createInvoiceQueue(additionalStackProps);
  }

  createProcessInvoiceQueue(additionalStackProps?: AdditionalStackProps) {
    return new Queue(this, "ProcessInvoice", {
      consumer: {
        function: {
          handler: "src/consumers/invoice_processor/main.go",
          environment: {
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            typlessToken: process.env.TYPLESS_TOKEN as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
          ],
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });
  }

  createInvoiceQueue(additionalStackProps?: AdditionalStackProps) {
    return new Queue(this, "Invoice", {
      consumer: {
        function: {
          handler: "src/consumers/invoice_preprocessor/main.go",
          environment: {
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            processInvoiceQueueUrl: this.processInvoiceQueue.sqsQueue
              .queueUrl as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
            this.processInvoiceQueue as Queue,
          ],
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });
  }
}
