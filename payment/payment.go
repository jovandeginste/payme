package payment

import (
	"fmt"
	"strings"

	"github.com/almerlucke/go-iban/iban"
)

// See: https://www.europeanpaymentscouncil.eu/document-library/guidance-documents/quick-response-code-guidelines-enable-data-capture-initiation

type Payment struct {
	ServiceTag string
	Version    int
	/*
		1: UTF-8 5: ISO 8859-5
		2: ISO 8859-1 6: ISO 8859-7
		3: ISO 8859-2 7: ISO 8859-10
		4: ISO 8859-4 8: ISO 8859-15
	*/
	CharacterSet       int
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

	RemittanceIsStructured bool
}

func NewStructured() Payment {
	return Payment{
		ServiceTag:             "BCD",
		Version:                2,
		CharacterSet:           2,
		IdentificationCode:     "SCT",
		RemittanceIsStructured: true,
	}
}

func New() Payment {
	return Payment{
		ServiceTag:             "BCD",
		Version:                2,
		CharacterSet:           2,
		IdentificationCode:     "SCT",
		RemittanceIsStructured: false,
	}
}

func (p *Payment) IBANBeneficiaryString() string {
	i, err := p.IBAN()
	if err != nil {
		return ""
	}

	return i.PrintCode
}

func (p *Payment) IBAN() (*iban.IBAN, error) {
	return iban.NewIBAN(p.IBANBeneficiary)
}

func (p *Payment) PurposeString() string {
	return strings.ReplaceAll(p.Purpose, " ", "")
}

func (p *Payment) VersionString() string {
	return fmt.Sprintf("%03d", p.Version)
}

func (p *Payment) CharacterSetString() string {
	return fmt.Sprintf("%d", p.Version)
}

func (p *Payment) EuroAmountString() string {
	if p.EuroAmount == 0 {
		return ""
	}

	return fmt.Sprintf("EUR%.2f", p.EuroAmount)
}

func (p *Payment) RemittanceStructured() string {
	return p.RemittanceString(true)
}

func (p *Payment) RemittanceText() string {
	return p.RemittanceString(false)
}

func (p *Payment) RemittanceString(structured bool) string {
	if p.RemittanceIsStructured != structured {
		return ""
	}

	return p.Remittance
}
