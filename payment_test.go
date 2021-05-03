package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncompletePayment(t *testing.T) {
	p := NewPayment()

	_, err := p.ToQRString()
	assert.Error(t, err)
}

func TestUnstructuredPayment(t *testing.T) {
	p := NewPayment()

	assert.Equal(t, "002", p.VersionString())
	assert.Equal(t, "2", p.CharacterSetString())

	p.NameBeneficiary = "François D'Alsace S.A."
	p.IBANBeneficiary = "FR1420041010050500013M02606"
	p.EuroAmount = 12.3
	p.Remittance = "Client:Marie Louise La Lune"

	result, err := p.ToQRString()
	assert.NoError(t, err)

	expected := `BCD
002
2
SCT

François D'Alsace S.A.
FR1420041010050500013M02606
EUR12.30


Client:Marie Louise La Lune
`

	assert.Equal(t, expected, result)
}

func TestStructuredPayment(t *testing.T) {
	p := NewStructuredPayment()

	p.Version = 1
	p.CharacterSet = 1
	p.BICBeneficiary = "BHBLDEHHXXX"
	p.NameBeneficiary = "Franz Mustermänn"
	p.IBANBeneficiary = "DE71110220330123456789"
	p.EuroAmount = 12.3
	p.Purpose = "GDDS"
	p.Remittance = "RF18539007547034"

	result, err := p.ToQRString()
	assert.NoError(t, err)

	expected := `BCD
001
1
SCT
BHBLDEHHXXX
Franz Mustermänn
DE71110220330123456789
EUR12.30
GDDS
RF18539007547034

`

	assert.Equal(t, expected, result)
}
