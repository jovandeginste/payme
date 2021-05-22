package payment

import (
	"errors"
	"regexp"
)

const (
	specialChars = `@&+()"':?.,-`
)

var StringValidator = regexp.MustCompile(`^[[:alnum:] ` + specialChars + `]+$`)

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

	if p.EuroAmount < 0.01 || p.EuroAmount > 999999999.99 {
		return errors.New("field 'EuroAmount' must be 0.01 or more and 999999999.99 or less")
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

	if !p.RemittanceIsStructured {
		if len(p.Remittance) > 140 {
			return errors.New("unstructured 'Remittance' should not exceed 140 characters")
		}

		if !StringValidator.MatchString(p.Remittance) {
			return errors.New("unstructured 'Remittance' should only contain alpha-numerics, spaces and/or " + specialChars)
		}
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

	if !StringValidator.MatchString(p.NameBeneficiary) {
		return errors.New("field 'NameBeneficiary' should not only contain alpha-numerics, spaces and/or " + specialChars)
	}

	if err := p.ValidateIBAN(); err != nil {
		return err
	}

	return nil
}

func (p *Payment) ValidateIBAN() error {
	_, err := p.IBAN()
	if err != nil {
		return err
	}

	return nil
}
