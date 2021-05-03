package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
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

func NewStructuredPayment() Payment {
	return Payment{
		ServiceTag:             "BCD",
		Version:                2,
		CharacterSet:           2,
		IdentificationCode:     "SCT",
		RemittanceIsStructured: true,
	}
}

func NewPayment() Payment {
	return Payment{
		ServiceTag:             "BCD",
		Version:                2,
		CharacterSet:           2,
		IdentificationCode:     "SCT",
		RemittanceIsStructured: false,
	}
}

func (p *Payment) ValidateHeader() error {
	if p.ServiceTag != "BCD" {
		return errors.New("field 'ServiceTag' should be BCD")
	}

	if p.CharacterSet < 1 || p.CharacterSet > 8 {
		return errors.New("field 'CharacterSet' should be 1..8")
	}

	if p.Version != 1 && p.Version != 2 {
		return errors.New("field 'Version' should be 1 or 2")
	}

	if p.IdentificationCode != "SCT" {
		return errors.New("field 'IdentificationCode' should be SCT")
	}

	if p.Version == 1 && p.BICBeneficiary == "" {
		return errors.New("field 'BICBeneficiary' is required when version is 1")
	}

	return nil
}

func (p *Payment) ValidateFields() error {
	if err := p.ValidateHeader(); err != nil {
		return err
	}

	if err := p.ValidateBeneficiary(); err != nil {
		return err
	}

	if p.EuroAmount != 0 {
		if p.EuroAmount < 0.01 || p.EuroAmount > 999999999.99 {
			return errors.New("field 'EuroAmount' must be 0.01 or more and 999999999.99 or less")
		}
	}

	if len(p.PurposeString()) > 4 {
		return errors.New("field 'Purpose' should not exceed 4 characters")
	}

	if err := p.ValidateRemittance(); err != nil {
		return err
	}

	return nil
}

func (p *Payment) ValidateRemittance() error {
	if p.Remittance == "" {
		return errors.New("field 'Remittance' is required")
	}

	if p.RemittanceIsStructured && len(p.Remittance) > 35 {
		return errors.New("structured 'Remittance' should not exceed 35 characters")
	}

	if !p.RemittanceIsStructured && len(p.Remittance) > 140 {
		return errors.New("unstructured 'Remittance' should not exceed 140 characters")
	}

	return nil
}

func (p *Payment) ValidateBeneficiary() error {
	if p.NameBeneficiary == "" {
		return errors.New("field 'NameBeneficiary' is required")
	}

	if len(p.NameBeneficiary) > 70 {
		return errors.New("field 'NameBeneficiary' should not exceed 70 characers")
	}

	if p.IBANBeneficiaryString() == "" {
		return errors.New("field 'IBANBeneficiary' is required")
	}

	if len(p.IBANBeneficiaryString()) > 34 {
		return errors.New("field 'IBANBeneficiary' should not exceed 70 characers")
	}

	return nil
}

func (p *Payment) ToQRString() (string, error) {
	if err := p.ValidateFields(); err != nil {
		return "", err
	}

	fields := []string{
		p.ServiceTag,
		p.VersionString(),
		p.CharacterSetString(),
		p.IdentificationCode,
		p.BICBeneficiary,
		p.NameBeneficiary,
		p.IBANBeneficiaryString(),
		p.EuroAmountString(),
		p.PurposeString(),
		p.RemittanceStructured(),
		p.RemittanceText(),
		p.B2OInformation,
	}

	return strings.Join(fields, "\n"), nil
}

func (p *Payment) IBANBeneficiaryString() string {
	return strings.ReplaceAll(p.IBANBeneficiary, " ", "")
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

func (p *Payment) ToQRPNG(qrSize int) ([]byte, error) {
	t, err := p.ToQRString()
	if err != nil {
		return nil, err
	}

	// Create the barcode
	qrCode, err := qr.Encode(t, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	// Scale the barcode to qrSize x qrSize pixels
	qrCode, err = barcode.Scale(qrCode, qrSize, qrSize)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	// encode the barcode as png
	err = png.Encode(&b, qrCode)

	return b.Bytes(), err
}
