package payment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	str40chars  = "1234567890123456789012345678901234567890"
	ExampleIBAN = "FR1420041010050500013M02606"
	ExampleName = "François D'Alsace S.A."
)

var (
	validStrings = []string{
		"My Name",
		ExampleName,
		"Franz Mustermänn",
		"The numb3r 0ne1",
		"Лев Николаевич Толстой",
		"老子",
		"이용철",
		"ألفت الادلبي",
		"M&M",
		"me@example.net",
		"How much is (5 + 1) : 3?",
		"One, 2, thr33",
		`"To be, or not to be."`,
	}
	invalidStrings = []string{
		"Don't allow {}",
		"No # symbol",
		"; drop table users",
		"🌈🌞🦄",
	}
)

func TestStringValidator(t *testing.T) {
	for _, s := range validStrings {
		assert.True(t, stringValidator.MatchString(s), fmt.Sprintf("Should match: %v", s))
	}

	for _, s := range invalidStrings {
		assert.False(t, stringValidator.MatchString(s), fmt.Sprintf("Should not match: %v", s))
	}
}

func validPayment() *Payment {
	p := New()
	p.Remittance = "Valid"
	p.NameBeneficiary = ExampleName
	p.IBANBeneficiary = ExampleIBAN

	return p
}

func TestValidateNew(t *testing.T) {
	p := validPayment()
	assert.NoError(t, p.validateHeader())

	p = NewStructured()
	assert.NoError(t, p.validateHeader())
}

func TestValidateHeader(t *testing.T) {
	p := validPayment()
	p.ServiceTag = "ABC"
	assert.ErrorIs(t, ErrValidationServiceTag, p.validateHeader())
	assert.ErrorIs(t, ErrValidationServiceTag, p.validateFields())

	p = validPayment()
	p.CharacterSet = 0
	assert.ErrorIs(t, ErrValidationCharacterSet, p.validateHeader())

	p = validPayment()
	p.CharacterSet = 9
	assert.ErrorIs(t, ErrValidationCharacterSet, p.validateHeader())

	p = validPayment()
	p.Version = 0
	assert.ErrorIs(t, ErrValidationVersion, p.validateHeader())

	p = validPayment()
	p.Version = 1
	p.BICBeneficiary = "XYZ"
	assert.NoError(t, p.validateHeader())

	p = validPayment()
	p.Version = 1
	p.BICBeneficiary = ""
	assert.ErrorIs(t, ErrValidationBICBeneficiary, p.validateHeader())

	p = validPayment()
	p.Version = 2
	assert.NoError(t, p.validateHeader())

	p = validPayment()
	p.Version = 3
	assert.ErrorIs(t, ErrValidationVersion, p.validateHeader())

	p = validPayment()
	p.IdentificationCode = "DEF"
	assert.ErrorIs(t, ErrValidationIdentificationCode, p.validateHeader())
}

func TestValidateFields(t *testing.T) {
	p := validPayment()

	for _, a := range []float64{-1, 0, 0.001, 0.00999, 999999999.991, 1000000000} {
		p.EuroAmount = a
		assert.ErrorIs(t, ErrValidationEuroAmount, p.validateFields(), fmt.Sprintf("Amount: %f", a))
	}

	for _, a := range []float64{0.01, 0.1, 1, 2.05, 99, 123456.78, 999999999.99} {
		p.EuroAmount = a
		assert.NoError(t, p.validateFields(), fmt.Sprintf("Amount: %f", a))
	}

	p = validPayment()
	p.EuroAmount = 1

	for _, n := range []string{"ABCDEF", "AB CD EF"} {
		p.Purpose = n
		assert.ErrorIs(t, ErrValidationPurpose, p.validateFields(), "Purpose: "+n)
	}

	for _, n := range []string{"ABCD", "AB CD", "A B C D"} {
		p.Purpose = n
		assert.NoError(t, p.validateFields())
	}
}

func TestValidateRemittance(t *testing.T) {
	p := validPayment()
	p.Remittance = ""
	p.EuroAmount = 1

	assert.ErrorIs(t, ErrValidationRemittanceRequired, p.validateRemittance())
	assert.ErrorIs(t, ErrValidationRemittanceRequired, p.validateFields())

	p.RemittanceIsStructured = true
	p.Remittance = str40chars
	assert.ErrorIs(t, ErrValidationRemittanceStructuredTooLong, p.validateRemittance())

	p.RemittanceIsStructured = false
	p.Remittance = str40chars + str40chars + str40chars + str40chars // 160 characters
	assert.ErrorIs(t, ErrValidationRemittanceUnstructuredTooLong, p.validateRemittance())

	p.RemittanceIsStructured = false
	p.Remittance = "#!"
	assert.ErrorIs(t, ErrValidationRemittanceUnstructuredCharacters, p.validateRemittance())
}

func TestValidateBeneficiary(t *testing.T) {
	p := validPayment()
	p.NameBeneficiary = ""
	assert.ErrorIs(t, ErrValidationNameBeneficiaryRequired, p.validateBeneficiary())

	p.NameBeneficiary = str40chars + str40chars // 80 characters
	assert.ErrorIs(t, ErrValidationNameBeneficiaryTooLong, p.validateBeneficiary())

	p.NameBeneficiary = "#!"
	assert.ErrorIs(t, ErrValidationNameBeneficiaryCharacters, p.validateBeneficiary())
}

func TestValidateIBAN(t *testing.T) {
	p := validPayment()
	p.IBANBeneficiary = "ABC"

	assert.Error(t, p.validateIBAN())
	assert.Error(t, p.validateBeneficiary())
}
