import {
  Stack,
  App,
  StackProps,
  Api,
  Queue,
  Table,
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
              jwtKey: process.env.JWT_KEY as string,
              tableName: additionalStackProps?.storageStack.invoiceTable
                .tableName as string,
            },
            permissions: [
              additionalStackProps?.storageStack.invoiceTable as Table,
            ],
          },
        },
        "POST /api/invoices/status": {
          function: {
            handler: "src/api/invoices/update_invoice_status/main.go",
            environment: {
              jwtKey: process.env.JWT_KEY as string,
              tableName: additionalStackProps?.storageStack.invoiceTable
                .tableName as string,
            },
            permissions: [
              additionalStackProps?.storageStack.invoiceTable as Table,
            ],
          },
        },
        "POST /api/invoices/data": {
          function: {
            handler: "src/api/invoices/update_invoice_data/main.go",
            environment: {
              jwtKey: process.env.JWT_KEY as string,
              tableName: additionalStackProps?.storageStack.invoiceTable
                .tableName as string,
            },
            permissions: [
              additionalStackProps?.storageStack.invoiceTable as Table,
            ],
          },
        },
      },
    });

    this.addOutputs({
      ApiEndpoint: api.url,
    });
  }
}
