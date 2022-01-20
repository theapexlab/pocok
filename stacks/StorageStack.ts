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
        // Keys
        pk: TableFieldType.STRING,
        sk: TableFieldType.STRING,
        lsi1sk: TableFieldType.STRING,
        lsi2sk: TableFieldType.STRING,
      },
      primaryIndex: { partitionKey: "pk", sortKey: "sk" },
      localIndexes: {
        invoiceStatusIndex: { sortKey: "lsi1sk" },
        customerEmailIndex: { sortKey: "lsi2sk" },
      },
    });
    
    this.invoiceBucket = new Bucket(this, "InvoiceBucket");
  }
}
