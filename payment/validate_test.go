package payment_test

import (
	"fmt"
	"testing"

	"github.com/jovandeginste/payme/payment"
	"github.com/stretchr/testify/assert"
)

const (
	str40chars = "1234567890123456789012345678901234567890"
)

var (
	validStrings = []string{
		"My Name",
		"François D'Alsace S.A.",
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
		assert.True(t, payment.StringValidator.MatchString(s), fmt.Sprintf("Should match: %v", s))
	}

	for _, s := range invalidStrings {
		assert.False(t, payment.StringValidator.MatchString(s), fmt.Sprintf("Should not match: %v", s))
	}
}

func validPayment() *payment.Payment {
	p := payment.New()
	p.Remittance = "Valid"
	p.NameBeneficiary = "François D'Alsace S.A."
	p.IBANBeneficiary = "FR1420041010050500013M02606"

	return p
}

func TestValidateNew(t *testing.T) {
	p := validPayment()
	assert.NoError(t, p.ValidateHeader())

	p = payment.NewStructured()
	assert.NoError(t, p.ValidateHeader())
}

func TestValidateHeader(t *testing.T) {
	p := validPayment()
	p.ServiceTag = "ABC"
	assert.ErrorIs(t, payment.ErrValidationServiceTag, p.ValidateHeader())
	assert.ErrorIs(t, payment.ErrValidationServiceTag, p.ValidateFields())

	p = validPayment()
	p.CharacterSet = 0
	assert.ErrorIs(t, payment.ErrValidationCharacterSet, p.ValidateHeader())

	p = validPayment()
	p.CharacterSet = 9
	assert.ErrorIs(t, payment.ErrValidationCharacterSet, p.ValidateHeader())

	p = validPayment()
	p.Version = 0
	assert.ErrorIs(t, payment.ErrValidationVersion, p.ValidateHeader())

	p = validPayment()
	p.Version = 1
	p.BICBeneficiary = "XYZ"
	assert.NoError(t, p.ValidateHeader())

	p = validPayment()
	p.Version = 1
	p.BICBeneficiary = ""
	assert.ErrorIs(t, payment.ErrValidationBICBeneficiary, p.ValidateHeader())

	p = validPayment()
	p.Version = 2
	assert.NoError(t, p.ValidateHeader())

	p = validPayment()
	p.Version = 3
	assert.ErrorIs(t, payment.ErrValidationVersion, p.ValidateHeader())

	p = validPayment()
	p.IdentificationCode = "DEF"
	assert.ErrorIs(t, payment.ErrValidationIdentificationCode, p.ValidateHeader())
}

func TestValidateFields(t *testing.T) {
	p := validPayment()

	for _, a := range []float64{-1, 0.001, 0.00999, 999999999.991, 1000000000} {
		p.EuroAmount = a
		assert.ErrorIs(t, payment.ErrValidationEuroAmount, p.ValidateFields(), fmt.Sprintf("Amount: %f", a))
	}

	for _, a := range []float64{0, 0.01, 0.1, 1, 2.05, 99, 123456.78, 999999999.99} {
		p.EuroAmount = a
		assert.NoError(t, p.ValidateFields(), fmt.Sprintf("Amount: %f", a))
	}

	p = validPayment()

	for _, n := range []string{"ABCDEF", "AB CD EF"} {
		p.Purpose = n
		assert.ErrorIs(t, payment.ErrValidationPurpose, p.ValidateFields(), "Purpose: "+n)
	}

	for _, n := range []string{"ABCD", "AB CD", "A B C D"} {
		p.Purpose = n
		assert.NoError(t, p.ValidateFields())
	}
}

func TestValidateRemittance(t *testing.T) {
	p := validPayment()
	p.Remittance = ""
	assert.ErrorIs(t, payment.ErrValidationRemittanceRequired, p.ValidateRemittance())
	assert.ErrorIs(t, payment.ErrValidationRemittanceRequired, p.ValidateFields())

	p.RemittanceIsStructured = true
	p.Remittance = str40chars
	assert.ErrorIs(t, payment.ErrValidationRemittanceStructuredTooLong, p.ValidateRemittance())

	p.RemittanceIsStructured = false
	p.Remittance = str40chars + str40chars + str40chars + str40chars // 160 characters
	assert.ErrorIs(t, payment.ErrValidationRemittanceUnstructuredTooLong, p.ValidateRemittance())

	p.RemittanceIsStructured = false
	p.Remittance = "#!"
	assert.ErrorIs(t, payment.ErrValidationRemittanceUnstructuredCharacters, p.ValidateRemittance())
}

func TestValidateBeneficiary(t *testing.T) {
	p := validPayment()
	p.NameBeneficiary = ""
	assert.ErrorIs(t, payment.ErrValidationNameBeneficiaryRequired, p.ValidateBeneficiary())

	p.NameBeneficiary = str40chars + str40chars // 80 characters
	assert.ErrorIs(t, payment.ErrValidationNameBeneficiaryTooLong, p.ValidateBeneficiary())

	p.NameBeneficiary = "#!"
	assert.ErrorIs(t, payment.ErrValidationNameBeneficiaryCharacters, p.ValidateBeneficiary())
}

func TestValidateIBAN(t *testing.T) {
	p := validPayment()
	p.IBANBeneficiary = "ABC"

	assert.Error(t, p.ValidateIBAN())
	assert.Error(t, p.ValidateBeneficiary())
}
