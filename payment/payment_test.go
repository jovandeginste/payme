package payment_test

import (
	"testing"

	"github.com/jovandeginste/payme/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncompletePayment(t *testing.T) {
	p := payment.New()

	_, err := p.ToString()
	require.Error(t, err)
}

func TestUnstructuredPayment(t *testing.T) {
	p := payment.New()

	assert.Equal(t, "002", p.VersionString())
	assert.Equal(t, "2", p.CharacterSetString())

	p.NameBeneficiary = ExampleName
	p.IBANBeneficiary = ExampleIBAN
	p.EuroAmount = 12.3
	p.Remittance = ExampleRemittance

	result, err := p.ToString()
	require.NoError(t, err)

	expected := `BCD
002
2
SCT

François D'Alsace S.A.
FR14 2004 1010 0505 0001 3M02 606
EUR12.30


Client:Marie Louise La Lune
`

	assert.Equal(t, expected, result)
}

func TestStructuredPayment(t *testing.T) {
	p := payment.NewStructured()

	p.Version = 1
	p.CharacterSet = 1
	p.BICBeneficiary = "BHBLDEHHXXX"
	p.NameBeneficiary = "Franz Mustermänn"
	p.IBANBeneficiary = "DE71110220330123456789"
	p.EuroAmount = 12.3
	p.Purpose = "GDDS"
	p.Remittance = "RF18539007547034"

	result, err := p.ToString()
	require.NoError(t, err)

	expected := `BCD
001
1
SCT
BHBLDEHHXXX
Franz Mustermänn
DE71 1102 2033 0123 4567 89
EUR12.30
GDDS
RF18539007547034

`

	assert.Equal(t, expected, result)
}

func TestEuroAmountString(t *testing.T) {
	p := payment.Payment{}
	p.EuroAmount = 0.01
	assert.Equal(t, "EUR0.01", p.EuroAmountString())

	p.EuroAmount = 1000.0001
	assert.Equal(t, "EUR1000.00", p.EuroAmountString())
}

func TestIBANBeneficiaryString(t *testing.T) {
	p := payment.Payment{}

	for _, s := range []string{
		"DE71110220330123456789",
		"DE71 11 022 0 33 0123 45678 9",
		"   DE71110220330123456789",
		"DE71110220330123456789  ",
	} {
		p.IBANBeneficiary = s
		assert.Equal(t, "DE71 1102 2033 0123 4567 89", p.IBANBeneficiaryString(), "From: "+s)
	}

	p.IBANBeneficiary = "invalid"
	assert.Empty(t, p.IBANBeneficiaryString())
}
