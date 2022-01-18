package mock

import "pocok/src/utils/models"

var MockInvoice = models.Invoice{
	Id:            "id",
	Filename:      "filename",
	TextractData:  ":)",
	InvoiceNumber: "500000",
	CustomerName:  "Csipkés Zoltán",
	AccountNumber: "10001000-10001000-10001000",
	Iban:          "HU69119800810030005009212644",
	NetPrice:      10000,
	GrossPrice:    20000,
	Currency:      "huf",
	DueDate:       "2050.01.01.",
	Services: []models.Service{
		{
			Name:       "Kutya",
			Amount:     500,
			NetPrice:   10,
			GrossPrice: 20,
			Currency:   "huf",
			Tax:        100,
		},
		{
			Name:       "Cica",
			Amount:     1000,
			NetPrice:   5,
			GrossPrice: 10,
			Currency:   "huf",
			Tax:        100,
		},
	},
	CustomerEmail: "sinter@example.com",
	Status:        models.PENDING,
}

var MockInvoice2 = models.Invoice{
	Id:            "id2",
	Filename:      "filename2",
	TextractData:  ":(",
	InvoiceNumber: "1",
	CustomerName:  "Wojak",
	AccountNumber: "10001000-10001000-10001000",
	Iban:          "PL69119800810030005009212644",
	NetPrice:      1,
	GrossPrice:    2,
	Currency:      "zł",
	DueDate:       "2050.01.01.",
	Services:      []models.Service{},
	CustomerEmail: "wojak@example.com",
	Status:        models.ACCEPTED,
}

var Invoices = []models.Invoice{
	MockInvoice, MockInvoice2,
}
