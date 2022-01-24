package mocks

import "pocok/src/utils/models"

var MockInvoice = models.Invoice{
	Pk:     "PK",
	Sk:     "SK",
	Lsi1sk: "LSI1SK",
	Lsi2sk: "LSI2SK",

	InvoiceId:  "ID",
	EntityType: models.INVOICE,
	Status:     models.PENDING,
	ReceivedAt: "2022.01.18",
	Filename:   "filename1",

	VendorEmail: "wojak@example.com",

	InvoiceNumber: "500000",
	VendorName:    "Csipkés Zoltán",
	AccountNumber: "10001000-10001000-10001000",
	Iban:          "HU69119800810030005009212644",
	NetPrice:      "10000",
	GrossPrice:    "20000",
	VatRate:       "27%",
	VatAmount:     "2700",
	Currency:      "huf",
	DueDate:       "2050.01.01.",
	Services: []models.Service{
		{
			Name:         "Kutya",
			Amount:       "500",
			UnitNetPrice: "10",
			NetPrice:     "5000",
			GrossPrice:   "10000",
			Currency:     "huf",
			VatRate:      "27%",
			VatAmount:    "2700",
		},
		{
			Name:         "Cica",
			Amount:       "1000",
			UnitNetPrice: "5",
			NetPrice:     "5000",
			GrossPrice:   "10000",
			Currency:     "huf",
			VatRate:      "27%",
			VatAmount:    "2700",
		},
	},
}

var MockInvoice2 = models.Invoice{
	Pk:     "PK",
	Sk:     "SK",
	Lsi1sk: "LSI1SK",
	Lsi2sk: "LSI2SK",

	InvoiceId:  "ID",
	EntityType: models.INVOICE,
	Status:     models.ACCEPTED,
	ReceivedAt: "2022.01.19",
	Filename:   "filename2",

	VendorEmail: "wojak@example.com",

	VatRate:       "AAM",
	VatAmount:     "0",
	InvoiceNumber: "1",
	VendorName:    "Wojak",
	AccountNumber: "10001000-10001000-10001000",
	Iban:          "PL69119800810030005009212644",
	NetPrice:      "1",
	GrossPrice:    "2",
	Currency:      "zł",
	DueDate:       "2050.01.01.",
	Services:      []models.Service{},
}

var Invoices = []models.Invoice{
	MockInvoice, MockInvoice2,
}