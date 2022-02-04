package models

type Quote struct {
	ID             string  `json:"id,omitempty"`
	Profile        int     `json:"profile"`
	TargetAccount  int     `json:"targetAccount"`
	SourceCurrency string  `json:"sourceCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	TargetAmount   int     `json:"targetAmount,omitempty"`
	SourceAmount   float64 `json:"sourceAmount,omitempty"`
	PayOut         string  `json:"payOut"`
}
