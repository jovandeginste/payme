package payment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, p.validateHeader())

	p = NewStructured()
	require.NoError(t, p.validateHeader())
}

func TestValidateHeader(t *testing.T) {
	p := validPayment()
	p.ServiceTag = "ABC"
	require.ErrorIs(t, ErrValidationServiceTag, p.validateHeader())
	require.ErrorIs(t, ErrValidationServiceTag, p.validateFields())

	p = validPayment()
	p.CharacterSet = 0
	require.ErrorIs(t, ErrValidationCharacterSet, p.validateHeader())

	p = validPayment()
	p.CharacterSet = 9
	require.ErrorIs(t, ErrValidationCharacterSet, p.validateHeader())

	p = validPayment()
	p.Version = 0
	require.ErrorIs(t, ErrValidationVersion, p.validateHeader())

	p = validPayment()
	p.Version = 1
	p.BICBeneficiary = "XYZ"
	require.NoError(t, p.validateHeader())

	p = validPayment()
	p.Version = 1
	p.BICBeneficiary = ""
	require.ErrorIs(t, ErrValidationBICBeneficiary, p.validateHeader())

	p = validPayment()
	p.Version = 2
	require.NoError(t, p.validateHeader())

	p = validPayment()
	p.Version = 3
	require.ErrorIs(t, ErrValidationVersion, p.validateHeader())

	p = validPayment()
	p.IdentificationCode = "DEF"
	require.ErrorIs(t, ErrValidationIdentificationCode, p.validateHeader())
}

func TestValidateFields(t *testing.T) {
	p := validPayment()

	for _, a := range []float64{-1, 0, 0.001, 0.00999, 999999999.991, 1000000000} {
		p.EuroAmount = a
		require.ErrorIs(t, ErrValidationEuroAmount, p.validateFields(), fmt.Sprintf("Amount: %f", a))
	}

	for _, a := range []float64{0.01, 0.1, 1, 2.05, 99, 123456.78, 999999999.99} {
		p.EuroAmount = a
		require.NoError(t, p.validateFields(), fmt.Sprintf("Amount: %f", a))
	}

	p = validPayment()
	p.EuroAmount = 1

	for _, n := range []string{"ABCDEF", "AB CD EF"} {
		p.Purpose = n
		require.ErrorIs(t, ErrValidationPurpose, p.validateFields(), "Purpose: "+n)
	}

	for _, n := range []string{"ABCD", "AB CD", "A B C D"} {
		p.Purpose = n
		require.NoError(t, p.validateFields())
	}
}

func TestValidateRemittance(t *testing.T) {
	p := validPayment()
	p.Remittance = ""
	p.EuroAmount = 1

	require.ErrorIs(t, ErrValidationRemittanceRequired, p.validateRemittance())
	require.ErrorIs(t, ErrValidationRemittanceRequired, p.validateFields())

	p.RemittanceIsStructured = true
	p.Remittance = str40chars
	require.ErrorIs(t, ErrValidationRemittanceStructuredTooLong, p.validateRemittance())

	p.RemittanceIsStructured = false
	p.Remittance = str40chars + str40chars + str40chars + str40chars // 160 characters
	require.ErrorIs(t, ErrValidationRemittanceUnstructuredTooLong, p.validateRemittance())

	p.RemittanceIsStructured = false
	p.Remittance = "#!"
	require.ErrorIs(t, ErrValidationRemittanceUnstructuredCharacters, p.validateRemittance())
}

func TestValidateBeneficiary(t *testing.T) {
	p := validPayment()
	p.NameBeneficiary = ""
	require.ErrorIs(t, ErrValidationNameBeneficiaryRequired, p.validateBeneficiary())

	p.NameBeneficiary = str40chars + str40chars // 80 characters
	require.ErrorIs(t, ErrValidationNameBeneficiaryTooLong, p.validateBeneficiary())

	p.NameBeneficiary = "#!"
	require.ErrorIs(t, ErrValidationNameBeneficiaryCharacters, p.validateBeneficiary())
}

func TestValidateIBAN(t *testing.T) {
	p := validPayment()
	p.IBANBeneficiary = "ABC"

	require.Error(t, p.validateIBAN())
	require.Error(t, p.validateBeneficiary())
}
