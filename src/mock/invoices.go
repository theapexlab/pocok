package mock

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

	CustomerEmail: "wojak@example.com",

	InvoiceNumber: "500000",
	CustomerName:  "Csipkés Zoltán",
	AccountNumber: "10001000-10001000-10001000",
	Iban:          "HU69119800810030005009212644",
	NetPrice:      10000,
	GrossPrice:    20000,
	Tax:           27,
	Currency:      "huf",
	DueDate:       "2050.01.01.",
	Services: []models.Service{
		{
			Name:         "Kutya",
			Amount:       500,
			UnitNetPrice: 10,
			NetPrice:     5000,
			GrossPrice:   10000,
			Currency:     "huf",
			Tax:          100,
		},
		{
			Name:         "Cica",
			Amount:       1000,
			UnitNetPrice: 5,
			NetPrice:     5000,
			GrossPrice:   10000,
			Currency:     "huf",
			Tax:          100,
		},
	},
	TextractData: ":)",
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

	CustomerEmail: "wojak@example.com",

	Tax:           27,
	InvoiceNumber: "1",
	CustomerName:  "Wojak",
	AccountNumber: "10001000-10001000-10001000",
	Iban:          "PL69119800810030005009212644",
	NetPrice:      1,
	GrossPrice:    2,
	Currency:      "zł",
	DueDate:       "2050.01.01.",
	Services:      []models.Service{},
	TextractData:  ":(",
}

var Invoices = []models.Invoice{
	MockInvoice, MockInvoice2,
}
