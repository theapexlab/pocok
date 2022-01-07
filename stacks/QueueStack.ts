import { Construct } from "@aws-cdk/core";
import { Queue, Stack, StackProps } from "@serverless-stack/resources";

export class QueueStack extends Stack {
  queue: Queue;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    this.queue = new Queue(this, "Queue");
  }
}
