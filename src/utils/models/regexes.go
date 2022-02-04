package models

const (
	GUESS_HUN_BANK_ACC    = "[0-9]{8}[ -]{1}[0-9]{8}([ -]{1}[0-9]{8})?$"
	VALIDATE_HUN_BANK_ACC = `^[0-9]{8}[\D]*[0-9]{8}([\D]*[0-9]{8})?$`
)
