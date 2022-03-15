package payment

import (
	"errors"
	"regexp"
)

const (
	specialChars = `@&+()"':?.,-/`
)

var (
	stringValidator = regexp.MustCompile(`^[\p{L}\d ` + specialChars + `]+$`)

	// ErrValidationServiceTag is returned when ServiceTag is not the correct value
	ErrValidationServiceTag = errors.New("field 'ServiceTag' should be BCD")
	// ErrValidationCharacterSet is returned when CharacterSet is not in the allowed range
	ErrValidationCharacterSet = errors.New("field 'CharacterSet' should be 1..8")
	// ErrValidationVersion is returned when Version is not 1 or 2
	ErrValidationVersion = errors.New("field 'Version' should be 1 or 2")
	// ErrValidationIdentificationCode is returned when IdentificationCode is not the correct value
	ErrValidationIdentificationCode = errors.New("field 'IdentificationCode' should be SCT")
	// ErrValidationBICBeneficiary is returned when BICBeneficiary is not set
	ErrValidationBICBeneficiary = errors.New("field 'BICBeneficiary' is required when version is 1")
	// ErrValidationEuroAmount is returned when EuroAmount is not a valid amount
	ErrValidationEuroAmount = errors.New("field 'EuroAmount' must be 0.01 or more and 999999999.99 or less")
	// ErrValidationPurpose is returned when Purpose is not within bounds
	ErrValidationPurpose = errors.New("field 'Purpose' should not exceed 4 characters")

	// ErrValidationRemittanceRequired is returned when Remittance is empty
	ErrValidationRemittanceRequired = errors.New("field 'Remittance' is required")
	// ErrValidationRemittanceStructuredTooLong is returned when Remittance is not within bounds for structured field
	ErrValidationRemittanceStructuredTooLong = errors.New("structured 'Remittance' should not exceed 35 characters")
	// ErrValidationRemittanceUnstructuredTooLong is returned when Remittance is not within bounds for unstructured field
	ErrValidationRemittanceUnstructuredTooLong = errors.New("unstructured 'Remittance' should not exceed 140 characters")
	// ErrValidationRemittanceUnstructuredCharacters is returned when Remittance contains invalid characters
	ErrValidationRemittanceUnstructuredCharacters = errors.New("unstructured 'Remittance' should only contain alpha-numerics, spaces and/or " + specialChars)

	// ErrValidationNameBeneficiaryRequired is returned when NameBeneficiary is empty
	ErrValidationNameBeneficiaryRequired = errors.New("field 'NameBeneficiary' is required")
	// ErrValidationNameBeneficiaryTooLong is returned when NameBeneficiary is not within bounds
	ErrValidationNameBeneficiaryTooLong = errors.New("field 'NameBeneficiary' should not exceed 70 characers")
	// ErrValidationNameBeneficiaryCharacters is returned when NameBeneficiary contains invalid characters
	ErrValidationNameBeneficiaryCharacters = errors.New("field 'NameBeneficiary' should not only contain alpha-numerics, spaces and/or " + specialChars)
)

// IsValid checks if all fields in the payment are consistent and meet the requirements.
// It returns the first error it encounters, or nil if all is well.
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

	if p.EuroAmount < 0.01 || p.EuroAmount > 999999999.99 {
		return ErrValidationEuroAmount
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

		if !stringValidator.MatchString(p.Remittance) {
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

	if !stringValidator.MatchString(p.NameBeneficiary) {
		return ErrValidationNameBeneficiaryCharacters
	}

	if err := p.validateIBAN(); err != nil {
		return err
	}

	return nil
}

func (p *Payment) validateIBAN() error {
	_, err := p.IBAN()
	return err
}
