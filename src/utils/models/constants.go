package models

// Organization ID
const (
	APEX_ID string = "ApexLab"
)

// InvoiceStatus Options
const (
	PENDING  string = "pending"
	ACCEPTED string = "approved"
)

// Field
const (
	STATUS   string = "STATUS"
	CUSTOMER string = "CUSTOMER"
	DATE     string = "DATE"
)

// Indexes
const (
	INVOICE_STATUS_INDEX string = "invoiceStatusIndex"
	CUSTOMER_EMAIL_INDEX string = "customerEmailIndex"
)

// Entity Types
const (
	ORG     string = "ORG"
	INVOICE string = "INVOICE"
)
