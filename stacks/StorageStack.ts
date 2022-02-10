import { RemovalPolicy } from "@aws-cdk/core";
import {
  App,
  Stack,
  Table,
  StackProps,
  TableFieldType,
  Bucket,
} from "@serverless-stack/resources";

export class StorageStack extends Stack {
  invoiceTable: Table;
  invoiceBucket: Bucket;
  assetBucket: Bucket;

  constructor(scope: App, id: string, props?: StackProps) {
    super(scope, id, props);

    this.invoiceTable = new Table(this, "Invoices", {
      fields: {
        // Keys
        pk: TableFieldType.STRING,
        sk: TableFieldType.STRING,
        lsi1sk: TableFieldType.STRING,
      },
      primaryIndex: { partitionKey: "pk", sortKey: "sk" },
      localIndexes: {
        localSecondaryIndex1: { sortKey: "lsi1sk" },
      },
      dynamodbTable: {
        removalPolicy: scope.local
          ? RemovalPolicy.DESTROY
          : RemovalPolicy.RETAIN,
      },
    });

    this.assetBucket = new Bucket(this, "AssetBucket");

    this.invoiceBucket = new Bucket(this, "InvoiceBucket", {
      s3Bucket: {
        removalPolicy: scope.local
          ? RemovalPolicy.DESTROY
          : RemovalPolicy.RETAIN,
      },
    });
  }
}
