import {
  Stack,
  App,
  StackProps,
  Api,
  Queue,
  Table,
  Bucket,
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

    const ampSharedEnvs = {
      jwtKey: process.env.JWT_KEY as string,
      mailgunSender: process.env.MAILGUN_SENDER as string,
    };

    const api = new Api(this, "Api", {
      customDomain:
        process.env.NODE_ENV === "development"
          ? undefined
          : {
              domainName: "api.pocok.biz",
              hostedZone: "pocok.biz",
            },
      routes: {
        "POST /webhooks/pipedream": {
          function: {
            handler: "src/api/process_email/main.go",
            environment: {
              queueUrl: additionalStackProps?.queueStack.invoiceQueue.sqsQueue
                .queueUrl as string,
            },
            permissions: [
              additionalStackProps?.queueStack.invoiceQueue as Queue,
            ],
          },
        },
        "GET /api/invoices": {
          function: {
            handler: "src/api/invoices/get_invoices/main.go",
            environment: {
              ...ampSharedEnvs,
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
        },
        "POST /api/invoices": {
          function: {
            handler: "src/api/invoices/update_invoice/main.go",
            environment: {
              ...ampSharedEnvs,
              tableName: additionalStackProps?.storageStack.invoiceTable
                .tableName as string,
              bucketName: additionalStackProps?.storageStack.invoiceBucket
                .bucketName as string,
              typlessToken: process.env.TYPLESS_TOKEN as string,
              typlessDocType: process.env.TYPLESS_DOC_TYPE as string,
            },
            permissions: [
              additionalStackProps?.storageStack.invoiceTable as Table,
              additionalStackProps?.storageStack.invoiceBucket as Bucket,
            ],
          },
        },
        "POST /api/invoices/accept_all": {
          function: {
            handler: "src/api/invoices/accept_all_invoices/main.go",
            environment: {
              ...ampSharedEnvs,
              tableName: additionalStackProps?.storageStack.invoiceTable
                .tableName as string,
            },
            permissions: [
              additionalStackProps?.storageStack.invoiceTable as Table,
            ],
          },
        },
        "POST /api/demo/invoice_summary": {
          function: {
            handler: "src/api/demo_cron/main.go",
            environment: {
              demoToken: process.env.DEMO_TOKEN as string,
              queueUrl: additionalStackProps?.queueStack.emailSenderQueue
                .sqsQueue.queueUrl as string,
            },
            permissions: [
              additionalStackProps?.queueStack.emailSenderQueue as Queue,
            ],
          },
        },
      },
      cors:
        process.env.NODE_ENV === "development"
          ? {
              allowOrigins: ["https://playground.amp.dev"],
            }
          : undefined,
    });

    this.addOutputs({
      ApiEndpoint: api.url,
    });
  }
}
