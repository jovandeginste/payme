package payment

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const QRSize = 300

func TestUnstructuredPaymentQR(t *testing.T) {
	p := New()

	assert.Equal(t, "002", p.VersionString())
	assert.Equal(t, "2", p.CharacterSetString())

	p.NameBeneficiary = ExampleName
	p.IBANBeneficiary = ExampleIBAN
	p.EuroAmount = 12.3
	p.Remittance = "Client:Marie Louise La Lune"

	result, err := p.ToQRString()
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile("tests/test1.qr")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	result, err = p.ToQRPNG(QRSize)
	assert.NoError(t, err)

	expected, err = ioutil.ReadFile("tests/test1.png")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStructuredPaymentQR(t *testing.T) {
	p := NewStructured()

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

	expected, err := ioutil.ReadFile("tests/test2.qr")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	result, err = p.ToQRPNG(QRSize)
	assert.NoError(t, err)

	expected, err = ioutil.ReadFile("tests/test2.png")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
