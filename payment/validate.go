package payment

import (
	"errors"
	"regexp"
)

const (
	specialChars = `@&+()"':?.,-`
)

var (
	StringValidator = regexp.MustCompile(`^[\p{L}\d ` + specialChars + `]+$`)

	ErrValidationServiceTag         = errors.New("field 'ServiceTag' should be BCD")
	ErrValidationCharacterSet       = errors.New("field 'CharacterSet' should be 1..8")
	ErrValidationVersion            = errors.New("field 'Version' should be 1 or 2")
	ErrValidationIdentificationCode = errors.New("field 'IdentificationCode' should be SCT")
	ErrValidationBICBeneficiary     = errors.New("field 'BICBeneficiary' is required when version is 1")
	ErrValidationEuroAmount         = errors.New("field 'EuroAmount' must be 0.01 or more and 999999999.99 or less")
	ErrValidationPurpose            = errors.New("field 'Purpose' should not exceed 4 characters")

	ErrValidationRemittanceRequired               = errors.New("field 'Remittance' is required")
	ErrValidationRemittanceStructuredTooLong      = errors.New("structured 'Remittance' should not exceed 35 characters")
	ErrValidationRemittanceUnstructuredTooLong    = errors.New("unstructured 'Remittance' should not exceed 140 characters")
	ErrValidationRemittanceUnstructuredCharacters = errors.New("unstructured 'Remittance' should only contain alpha-numerics, spaces and/or " + specialChars)

	ErrValidationNameBeneficiaryRequired   = errors.New("field 'NameBeneficiary' is required")
	ErrValidationNameBeneficiaryTooLong    = errors.New("field 'NameBeneficiary' should not exceed 70 characers")
	ErrValidationNameBeneficiaryCharacters = errors.New("field 'NameBeneficiary' should not only contain alpha-numerics, spaces and/or " + specialChars)
)

// Validate checks if all fields in the payment are consistent and meet the requirements
func (p *Payment) IsValid() error {
	return p.validateFields()
}

func (p *Payment) validateFields() error {
	if err := p.validateHeader(); err != nil {
		return err
	}

	if err := p.validateBeneficiary(); err != nil {
		return err
	}

	if p.EuroAmount != 0 {
		if p.EuroAmount < 0.01 || p.EuroAmount > 999999999.99 {
			return ErrValidationEuroAmount
		}
	}

	if len(p.PurposeString()) > 4 {
		return ErrValidationPurpose
	}

	if err := p.validateRemittance(); err != nil {
		return err
	}

	return nil
}

func (p *Payment) validateHeader() error {
	if p.ServiceTag != "BCD" {
		return ErrValidationServiceTag
	}

	if p.CharacterSet < 1 || p.CharacterSet > 8 {
		return ErrValidationCharacterSet
	}

	if p.Version != 1 && p.Version != 2 {
		return ErrValidationVersion
	}

	if p.IdentificationCode != "SCT" {
		return ErrValidationIdentificationCode
	}

	if p.Version == 1 && p.BICBeneficiary == "" {
		return ErrValidationBICBeneficiary
	}

	return nil
}

func (p *Payment) validateRemittance() error {
	if p.Remittance == "" {
		return ErrValidationRemittanceRequired
	}

	if p.RemittanceIsStructured && len(p.Remittance) > 35 {
		return ErrValidationRemittanceStructuredTooLong
	}

	if !p.RemittanceIsStructured {
		if len(p.Remittance) > 140 {
			return ErrValidationRemittanceUnstructuredTooLong
		}

		if !StringValidator.MatchString(p.Remittance) {
			return ErrValidationRemittanceUnstructuredCharacters
		}
	}

	return nil
}

func (p *Payment) validateBeneficiary() error {
	if p.NameBeneficiary == "" {
		return ErrValidationNameBeneficiaryRequired
	}

	if len(p.NameBeneficiary) > 70 {
		return ErrValidationNameBeneficiaryTooLong
	}

	if !StringValidator.MatchString(p.NameBeneficiary) {
		return ErrValidationNameBeneficiaryCharacters
	}

	if err := p.validateIBAN(); err != nil {
		return err
	}

	return nil
}

func (p *Payment) validateIBAN() error {
	_, err := p.IBAN()
	if err != nil {
		return err
	}

	return nil
}
