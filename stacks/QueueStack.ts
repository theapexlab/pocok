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
  invoiceQueue: Queue;
  processInvoiceQueue: Queue;
  emailSenderQueue: Queue;

  constructor(
    scope: Construct,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    this.processInvoiceQueue =
      this.createProcessInvoiceQueue(additionalStackProps);
    this.invoiceQueue = this.createPreprocessInvoiceQueue(additionalStackProps);
    this.emailSenderQueue = this.createEmailSenderQueue(additionalStackProps);
  }

  createProcessInvoiceQueue(additionalStackProps?: AdditionalStackProps) {
    return new Queue(this, "ProcessInvoice", {
      consumer: {
        function: {
          handler: "src/consumers/invoice_processor/main.go",
          environment: {
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            typlessToken: process.env.TYPLESS_TOKEN  || "",
            typlessDocType: process.env.TYPLESS_DOC_TYPE || "",
            tableName: additionalStackProps?.storageStack.invoiceTable
            .tableName as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
            additionalStackProps?.storageStack.invoiceTable as Table
          ],
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

  createEmailSenderQueue(additionalStackProps?: AdditionalStackProps) {
    return new Queue(this, "EmailSender", {
      consumer: {
        function: {
          handler: "src/consumers/email_sender/main.go",
          environment: {
            sender: process.env.MAILGUN_SENDER as string,
            mailgunDomain: process.env.MAILGUN_DOMAIN as string,
            mailgunApiKey: process.env.MAILGUN_API_KEY as string,
            emailRecipient: process.env.EMAIL_RECIPIENT as string,
            apiUrl: process.env.API_URL as string,
            jwtKey: process.env.JWT_KEY as string,
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
