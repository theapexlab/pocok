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
      customDomain: scope.local
        ? undefined
        : {
            domainName: process.env.DOMAIN_NAME as string,
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
        "POST /api/invoices/status": {
          function: {
            handler: "src/api/invoices/update_invoice_status/main.go",
            environment: {
              ...ampSharedEnvs,
              tableName: additionalStackProps?.storageStack.invoiceTable
                .tableName as string,
              typlessToken: process.env.TYPLESS_TOKEN as string,
              typlessDocType: process.env.TYPLESS_DOC_TYPE as string,
              bucketName: additionalStackProps?.storageStack.invoiceBucket
                .bucketName as string,
            },
            permissions: [
              additionalStackProps?.storageStack.invoiceTable as Table,
              additionalStackProps?.storageStack.invoiceBucket as Bucket,
            ],
          },
        },
        "POST /api/invoices/data": {
          function: {
            handler: "src/api/invoices/update_invoice_data/main.go",
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
      cors: scope.local
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
