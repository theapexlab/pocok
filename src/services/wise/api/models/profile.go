package models

import "fmt"

type Profile struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Details struct {
		Name               string `json:"name"`
		RegistrationNumber string `json:"registrationNumber"`
		CompanyType        string `json:"companyType"`
		CompanyRole        string `json:"companyRole"`
	} `json:"details"`
}

func (p Profile) String() string {
	return fmt.Sprintf("%v %v %v", p.ID, p.Type, p.Details.Name)
}
