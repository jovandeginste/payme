package payment_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jovandeginste/payme/payment"
	"github.com/stretchr/testify/assert"
)

const (
	QRSize            = 300
	ExampleIBAN       = "FR1420041010050500013M02606"
	ExampleName       = "François D'Alsace S.A."
	ExampleRemittance = "Client:Marie Louise La Lune"
)

func TestInvalid(t *testing.T) {
	p := payment.New()

	err := p.IsValid()
	assert.Error(t, err)

	_, err = p.ToQRString()
	assert.Error(t, err)

	_, err = p.ToQRPNG(QRSize)
	assert.Error(t, err)
}

func ExamplePayment() {
	p := payment.New()

	p.NameBeneficiary = ExampleName
	p.IBANBeneficiary = ExampleIBAN
	p.EuroAmount = 12.3
	p.Remittance = ExampleRemittance

	fmt.Println(p.ToQRString())
}

func TestUnstructuredPaymentQR(t *testing.T) {
	p := payment.New()

	assert.Equal(t, "002", p.VersionString())
	assert.Equal(t, "2", p.CharacterSetString())

	p.NameBeneficiary = ExampleName
	p.IBANBeneficiary = ExampleIBAN
	p.EuroAmount = 12.3
	p.Remittance = ExampleRemittance

	err := p.IsValid()
	assert.NoError(t, err)

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
	p := payment.NewStructured()

	p.Version = 1
	p.CharacterSet = 1
	p.BICBeneficiary = "BHBLDEHHXXX"
	p.NameBeneficiary = "Franz Mustermänn"
	p.IBANBeneficiary = "DE71110220330123456789"
	p.EuroAmount = 12.3
	p.Purpose = "GDDS"
	p.Remittance = "RF18539007547034"

	err := p.IsValid()
	assert.NoError(t, err)

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
