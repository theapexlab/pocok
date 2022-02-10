import { Duration } from "@aws-cdk/core";
import {
  App,
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
  invoiceQueue: Queue;
  processInvoiceQueue: Queue;
  emailSenderQueue: Queue;
  wiseQueue: Queue;

  constructor(
    scope: App,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    this.processInvoiceQueue =
      this.createProcessInvoiceQueue(additionalStackProps);
    this.invoiceQueue = this.createPreprocessInvoiceQueue(additionalStackProps);
    this.emailSenderQueue = this.createEmailSenderQueue({
      scope,
      additionalStackProps,
    });
    this.wiseQueue = this.createWiseQueue();
  }

  createProcessInvoiceQueue(additionalStackProps?: AdditionalStackProps) {
    const lambdaTimeout =
      process.env.PROCESS_INVOICE_LAMBDA_TIMEOUT_SEC || "60";
    return new Queue(this, "ProcessInvoice", {
      consumer: {
        function: {
          handler: "src/consumers/invoice_processor/main.go",
          environment: {
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            typlessToken: process.env.TYPLESS_TOKEN || "",
            typlessDocType: process.env.TYPLESS_DOC_TYPE || "",
            tableName: additionalStackProps?.storageStack.invoiceTable
              .tableName as string,
            lambdaTimeout,
            mailgunSender: process.env.MAILGUN_SENDER as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
            additionalStackProps?.storageStack.invoiceTable as Table,
          ],
          // FYI: default 6s is may not enough for typless requests to complete
          timeout: Duration.seconds(parseInt(lambdaTimeout)),
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });
  }

  createPreprocessInvoiceQueue(additionalStackProps?: AdditionalStackProps) {
    return new Queue(this, "Invoice", {
      consumer: {
        function: {
          handler: "src/consumers/invoice_preprocessor/main.go",
          environment: {
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            processInvoiceQueueUrl: this.processInvoiceQueue.sqsQueue
              .queueUrl as string,
            tableName: additionalStackProps?.storageStack.invoiceTable
              .tableName as string,
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

  createEmailSenderQueue({
    scope,
    additionalStackProps,
  }: {
    additionalStackProps?: AdditionalStackProps;
    scope: App;
  }) {
    return new Queue(this, "EmailSender", {
      consumer: {
        function: {
          handler: "src/consumers/email_sender/main.go",
          environment: {
            mailgunSender: process.env.MAILGUN_SENDER as string,
            mailgunDomain: process.env.MAILGUN_DOMAIN as string,
            mailgunApiKey: process.env.MAILGUN_API_KEY as string,
            emailRecipient: process.env.EMAIL_RECIPIENT as string,
            apiUrl: process.env.API_URL as string,
            jwtKey: process.env.JWT_KEY as string,
            assetBucketName: additionalStackProps?.storageStack.assetBucket
              .bucketName as string,
            stage: scope.local ? "development" : "production",
          },
          permissions: [
            additionalStackProps?.storageStack.assetBucket as Bucket,
          ],
          bundle: {
            copyFiles: [{ from: "src/amp/templates", to: "src/amp/templates" }],
          },
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });
  }

  createWiseQueue() {
    const wiseQueue = new Queue(this, "Wise");
    wiseQueue.addConsumer(this, {
      function: {
        handler: "src/consumers/wise_processor/main.go",
        environment: {
          queueUrl: wiseQueue.sqsQueue.queueUrl,
          wiseApiToken: process.env.WISE_API_TOKEN as string,
        },
        permissions: [wiseQueue],
      },
      consumerProps: {
        batchSize: 1,
      },
    });

    return wiseQueue;
  }
}
