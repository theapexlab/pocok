import { Construct } from "@aws-cdk/core";
import {
  Bucket,
  Queue,
  Stack,
  StackProps,
  Table,
  Topic
} from "@serverless-stack/resources";
import { StorageStack } from "./StorageStack";
import { PolicyStatement, Effect, Role, ServicePrincipal } from "@aws-cdk/aws-iam"

type AdditionalStackProps = {
  storageStack: StorageStack;
};

export class QueueStack extends Stack {
  invoiceQueue: Queue;
  textractJobResultsQueue: Queue;
  textractJobCompletionTopic: Topic;

  constructor(
    scope: Construct,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    this.textractJobResultsQueue = this.createTextractJobResultsQueue(additionalStackProps)
    this.textractJobCompletionTopic = this.createTextractJobCompletionTopic()
    this.invoiceQueue = this.createInvoiceQueue(additionalStackProps)
  }

  createTextractJobResultsQueue(additionalStackProps?: AdditionalStackProps) {
    return  new Queue(this, "TextractJobResults", {
      consumer: {
        function: {
          handler: "src/consumers/textractor/main.go",
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
          initialPolicy: [
            new PolicyStatement({
              resources: ["*"],
              actions: ["textract:*"]
            })
          ]
        },
        consumerProps: {
          batchSize: 1,
        },
      }
    });
  }

  createTextractJobCompletionTopic() {
    return new Topic(this, "TextractJobCompletion", {
      subscribers: [
        this.textractJobResultsQueue as Queue,
      ],
    });
  }

  createInvoiceQueue(additionalStackProps?: AdditionalStackProps) {
    const textractServiceRole = new Role(this, "textractServiceRole", {
      assumedBy: new ServicePrincipal("textract.amazonaws.com"),
    });

    textractServiceRole.addToPolicy(
      new PolicyStatement({
        effect: Effect.ALLOW,
        resources: [this.textractJobCompletionTopic.topicArn as string],
        actions: ["sns:Publish"],
      })
    );

    return new Queue(this, "Invoice", {
      consumer: {
        function: {
          handler: "src/consumers/invoice_preprocessor/main.go",
          environment: {
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            textractQueueUrl:  this.textractJobResultsQueue.sqsQueue.queueUrl  as string,
            snsTopicArn: this.textractJobCompletionTopic.topicArn as string,
            textractRoleArn: textractServiceRole.roleArn as string
          },
          permissions: [
            this.textractJobResultsQueue as Queue,
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
          ],
          initialPolicy: [
            new PolicyStatement({
              resources: ["*"],
              actions: ["textract:*"]
            })
          ]
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });
  }
}
