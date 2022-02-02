package models

type Transfer struct {
	ID                    int    `json:"id,omitempty"`
	TargetAccount         int    `json:"targetAccount"`
	QuoteUUID             string `json:"quoteUuid"`
	CustomerTransactionID string `json:"customerTransactionId"`
	Details               struct {
		Reference string `json:"reference"`
	} `json:"details"`
}
