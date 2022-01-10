import { Construct } from "@aws-cdk/core";
import {
  Stack,
  Table,
  StackProps,
  TableFieldType,
  Bucket,
} from "@serverless-stack/resources";

export class StorageStack extends Stack {
  invoiceTable: Table;
  invoiceBucket: Bucket;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    this.invoiceTable = new Table(this, "Invoices", {
      fields: {
        id: TableFieldType.STRING,
        filename: TableFieldType.STRING,
        etag: TableFieldType.STRING,
      },
      primaryIndex: { partitionKey: "id" },
    });

    this.invoiceBucket = new Bucket(this, "InvoiceBucket");
  }
}
