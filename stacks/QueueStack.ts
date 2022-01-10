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
            s3Bucket: additionalStackProps?.storageStack.invoiceBucket
              .bucketName as string,
          },
          permissions: [
            additionalStackProps?.storageStack.invoiceTable as Table,
            additionalStackProps?.storageStack.invoiceBucket as Bucket,
          ],
        },
      },
    });
  }
}
