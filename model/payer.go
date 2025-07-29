package model

import "time"

type Payer struct {
	Name              string    `json:"name"`                // Name of the tenant
	Phone             string    `json:"phone"`               // Tenant's phone number
	RentalEmail       string    `json:"rental_email"`        // Tenant's email for notifications
	RentalDate        time.Time `json:"rental_date"`         // Date of rent start
	RenterName        string    `json:"renter_name"`         // Name of the property owner/landlord
	RenterEmail       string    `json:"renter_email"`        // Email of the landlord
	NIT               string    `json:"nit"`                 // Tenant's tax ID or national ID
	PropertyAddress   string    `json:"property_address"`    // Full address of the rented property
	PropertyType      string    `json:"property_type"`       // Type of property (e.g., Apartment, House)
	RentalStart       time.Time `json:"rental_start"`        // Rental contract start date
	RentalEnd         time.Time `json:"rental_end"`          // Rental contract end date
	MonthlyRent       int       `json:"monthly_rent"`        // Monthly rental amount
	BankName          string    `json:"bank_name"`           // Name of the bank for payment
	AccountType       string    `json:"account_type"`        // Type of account (Savings/Current)
	BankAccountNumber string    `json:"bank_account_number"` // Bank account number for payments
	AccountHolder     string    `json:"account_holder"`      // Name of the bank account holder
	PaymentTerms      string    `json:"payment_terms"`       // Payment conditions/details
	AdditionalNotes   string    `json:"additional_notes"`    // Any extra observations
	UnpaidMonths      int       `json:"unpaid_months"`
}
