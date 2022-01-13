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
        textractData: TableFieldType.STRING,

        // Extracted Data
        invoiceNumber: TableFieldType.STRING,
        customerName: TableFieldType.STRING,
        accountNumber: TableFieldType.STRING,
        iban: TableFieldType.STRING,
        netPrice: TableFieldType.NUMBER,
        grossPrice: TableFieldType.NUMBER,
        currency: TableFieldType.STRING,
        dueDate: TableFieldType.STRING,

        // Refactor into service array/table later
        serviceName: TableFieldType.STRING,
        serviceAmount: TableFieldType.STRING,
        serviceNetPrice: TableFieldType.STRING,
        serviceGrossPrice: TableFieldType.STRING,
        serviceCurrency: TableFieldType.STRING,
        serviceTax: TableFieldType.STRING,

        // Dynamic data
        customerEmail: TableFieldType.STRING,
        status: TableFieldType.STRING,
      },
      primaryIndex: { partitionKey: "id", sortKey: "filename" },
      localIndexes: {
        invoiceNumberIndex: { sortKey: "invoiceNumber" },
        customerIndex: { sortKey: "customerName" },
      },
    });
    this.invoiceBucket = new Bucket(this, "InvoiceBucket");
  }
}
