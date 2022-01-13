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
  uploadQueue: Queue;
  textractQueue: Queue;
  jobCompletionTopic: Topic;

  constructor(
    scope: Construct,
    id: string,
    props?: StackProps,
    additionalStackProps?: AdditionalStackProps
  ) {
    super(scope, id, props);

    this.textractQueue = new Queue(this, "textractQueue", {
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
        },
        consumerProps: {
          batchSize: 1,
        },
      }
    });

    this.jobCompletionTopic = new Topic(this, "jobCompletionTopic", {
      subscribers: [
        this.textractQueue as Queue,
      ],
    });

    const textractServiceRole = new Role(this, "textractServiceRole", {
      assumedBy: new ServicePrincipal("textract.amazonaws.com"),
    });

    textractServiceRole.addToPolicy(
      new PolicyStatement({
        effect: Effect.ALLOW,
        resources: [this.jobCompletionTopic.topicArn as string],
        actions: ["sns:Publish"],
      })
    );

    // todo: remove un-neccesary environment vars
    this.uploadQueue = new Queue(this, "uploadQueue", {
      consumer: {
        function: {
          handler: "src/consumers/s3_uploader/main.go",
          environment: {
            tableName: additionalStackProps?.storageStack.invoiceTable
              .tableName as string,
            bucketName: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
            textractQueueUrl:  this.textractQueue.sqsQueue.queueUrl  as string,
            snsTopicArn: this.jobCompletionTopic.topicArn as string,
            textractRoleArn: textractServiceRole.roleArn as string
          },
          permissions: [
            this.textractQueue as Queue,
            additionalStackProps?.storageStack.invoiceTable as Table,
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
          ],
        },
        consumerProps: {
          batchSize: 1,
        },
      },
    });

    const textractRole = new Role(this, "TextractServiceRole", {
      assumedBy: new ServicePrincipal("textract.amazonaws.com"),
    })

    this.uploadQueue.consumerFunction?.addToRolePolicy(new PolicyStatement({
      resources: ["*"],
      actions: ["textract:*"]
    }))
  }
}
