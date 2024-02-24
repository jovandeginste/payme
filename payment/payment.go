package payment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/almerlucke/go-iban/iban"
)

// See: https://www.europeanpaymentscouncil.eu/document-library/guidance-documents/quick-response-code-guidelines-enable-data-capture-initiation

// Payment encapsulates all fields needed to generate the QR code
type Payment struct {
	// ServiceTag should always be BCD
	ServiceTag string
	// Version should be v1 or v2
	Version int
	/*
		1: UTF-8 5: ISO 8859-5
		2: ISO 8859-1 6: ISO 8859-7
		3: ISO 8859-2 7: ISO 8859-10
		4: ISO 8859-4 8: ISO 8859-15
	*/
	CharacterSet int
	// IdentificationCode should always be SCT (SEPA Credit Transfer)
	IdentificationCode string
	// AT-23 BIC of the Beneficiary Bank [optional in Version 2]
	// The BIC will continue to be mandatory for SEPA payment transactions involving non-EEA countries.
	BICBeneficiary string
	// AT-21 Name of the Beneficiary
	NameBeneficiary string
	// AT-20 Account number of the Beneficiary
	// Only IBAN is allowed.
	IBANBeneficiary string
	// AT-04 Amount of the Credit Transfer in Euro [optional]
	// Amount must be 0.01 or more and 999999999.99 or less
	EuroAmount float64
	// AT-44 Purpose of the Credit Transfer [optional]
	Purpose string
	// AT-05 Remittance Information (Structured) [optional]
	// Creditor Reference (ISO 11649 RFCreditor Reference may be used
	// *or*
	// AT-05 Remittance Information (Unstructured) [optional]
	Remittance string
	// Beneficiary to originator information [optional]
	B2OInformation string

	// Defines whether the Remittance Information is Structured or Unstructured
	RemittanceIsStructured bool
}

// NewStructured returns a default Payment with the Structured flag enabled
func NewStructured() *Payment {
	p := New()

	p.RemittanceIsStructured = true

	return p
}

// New returns a new Payment struct with default values for version 2
func New() *Payment {
	return &Payment{
		ServiceTag:             "BCD",
		Version:                2,
		CharacterSet:           2,
		IdentificationCode:     "SCT",
		RemittanceIsStructured: false,
	}
}

// IBANBeneficiaryString returns the IBAN of the beneficiary in a standardized form
func (p *Payment) IBANBeneficiaryString() string {
	i, err := p.IBAN()
	if err != nil {
		return ""
	}

	return i.PrintCode
}

// IBAN returns the parsed IBAN of the beneficiary
func (p *Payment) IBAN() (*iban.IBAN, error) {
	return iban.NewIBAN(p.IBANBeneficiary)
}

// PurposeString returns the parsed purpose
func (p *Payment) PurposeString() string {
	return strings.ReplaceAll(p.Purpose, " ", "")
}

// VersionString returns the version converted to a 3-digit number with leading zeros
func (p *Payment) VersionString() string {
	return fmt.Sprintf("%03d", p.Version)
}

// CharacterSetString returns the character set converted to string
func (p *Payment) CharacterSetString() string {
	return strconv.Itoa(p.CharacterSet)
}

// EuroAmountString returns the set amount in financial format (eg. EUR12.34)
// or an empty string if the amount is 0
func (p *Payment) EuroAmountString() string {
	return fmt.Sprintf("EUR%.2f", p.EuroAmount)
}

// RemittanceStructured returns the value for the structured remittance line
func (p *Payment) RemittanceStructured() string {
	return p.RemittanceString(true)
}

// RemittanceText returns the value for the unstructured (freeform) remittance line
func (p *Payment) RemittanceText() string {
	return p.RemittanceString(false)
}

// RemittanceString returns the value for the remittance field, independing on being structured
func (p *Payment) RemittanceString(structured bool) string {
	if p.RemittanceIsStructured != structured {
		return ""
	}

	return p.Remittance
}

// BICBeneficiaryString returns the BIC of the beneficiary, depending on the version of the QR code
func (p *Payment) BICBeneficiaryString() string {
	if p.Version != 1 {
		return ""
	}

	return p.BICBeneficiary
}
