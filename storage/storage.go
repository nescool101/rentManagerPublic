package storage

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nescool101/rentManager/model"
)

const filePath = "payers.json"

func InitializePayersFile() {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		data := []byte(`[
  {
    "name": "Angelica Dominguez",
    "phone": "57 300 4872915",
    "rental_email": "angedo2210@hotmail.com",
    "rental_date": "2022-06-15T00:00:00Z",
    "renter_name": "Nestor Fernando Alvarez Gomez",
    "renter_email": "nescool101@gmail.com",
    "nit": "1015398879",
    "property_address": "Cra. 56 #167-29 501",
    "property_type": "Apartmento",
    "rental_start": "2023-12-01T00:00:00Z",
    "rental_end": "2024-12-01T00:00:00Z",
    "monthly_rent": 1100000,
    "bank_name": "Banco Caja Social",
    "account_type": "Ahorros",
    "bank_account_number": "2405896900",
    "account_holder": "Nestor Fernando Alvarez Gomez",
    "payment_terms": "Pago se realiza en 2 cuotas, el 15 y el 1ro de cada mes.",
    "additional_notes": "Late payments incur a 5% penalty.",
    "unpaid_months": 0
  },
  {
    "name": "Stiwar Cortes",
    "phone": "57 322 8761776",
    "rental_email": "Stiwuar1011cortes@gmail.com",
    "rental_date": "2023-07-28T00:00:00Z",
    "renter_name": "Nestor Fernando Alvarez Gomez",
    "renter_email": "nescool101@gmail.com",
    "nit": "1015398879",
    "property_address": "Cl 167 #56 25 Apto 1 203",
    "property_type": "Apartmento",
    "rental_start": "2023-12-01T00:00:00Z",
    "rental_end": "2024-12-01T00:00:00Z",
    "monthly_rent": 1200000,
    "bank_name": "Banco Caja Social",
    "account_type": "Ahorros",
    "bank_account_number": "2405896900",
    "account_holder": "Nestor Fernando Alvarez Gomez",
    "payment_terms": "Pago debe ser en transferencia bancaria.",
    "additional_notes": "Late payments incur a 5% penalty.",
    "unpaid_months": 0
  },
  {
    "name": "testnombre test mes",
    "phone": "234-567-8901",
    "rental_email": "nescool101@gmail.com",
    "rental_date": "2023-01-07T00:00:00Z",
    "renter_name": "Carlos Owner test",
    "renter_email": "nescool101@gmail.com",
    "nit": "800234567",
    "property_address": "456 Oak St, Townsville",
    "property_type": "House",
    "rental_start": "2024-06-01T00:00:00Z",
    "rental_end": "2025-06-01T00:00:00Z",
    "monthly_rent": 1100000,
    "bank_name": "Chase Bank",
    "account_type": "Savings",
    "bank_account_number": "987654321",
    "account_holder": "Jane Smith",
    "payment_terms": "Payments must be made via bank transfer.",
    "additional_notes": "Tenant must notify 30 days before contract termination.",
    "unpaid_months": 0
  }
]`)
		err := os.WriteFile(filePath, data, 0644)
		if err != nil {
			log.Fatalf("Failed to create payers file: %v", err)
		}
	}
}

func GetPayers() ([]model.Payer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var payers []model.Payer
	err = json.Unmarshal(data, &payers)
	return payers, err
}

func parseDate(dateStr string) time.Time {
	layout := time.RFC3339
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Fatalf("Invalid date format for %s", dateStr)
	}
	return parsedDate
}
