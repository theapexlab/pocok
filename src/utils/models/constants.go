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
	STATUS string = "STATUS"
	VENDOR string = "VENDOR"
	DATE   string = "DATE"
)

// Indexes
const (
	LOCAL_SECONDARY_INDEX_1 string = "localSecondaryIndex1"
	LOCAL_SECONDARY_INDEX_2 string = "localSecondaryIndex2"
)

// Entity Types
const (
	ORG     string = "ORG"
	INVOICE string = "INVOICE"
)

// Email Event Types
const (
	EMAIL_SUMMARY string = "EMAIL_SUMMARY"
)

// Email Contents
const (
	EMAIL_SUMMARY_SUBJECT string = "Daily invoice summary email"
	EMAIL_NO_AMP_BODY     string = "This email requires AMP to be enabled"
)
