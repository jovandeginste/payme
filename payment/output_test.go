package payment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jovandeginste/payme/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.Error(t, err)

	_, err = p.ToQRBytes()
	require.Error(t, err)

	_, err = p.ToQRPNG(QRSize)
	require.Error(t, err)
}

func ExamplePayment() {
	p := payment.New()

	p.NameBeneficiary = ExampleName
	p.IBANBeneficiary = ExampleIBAN
	p.EuroAmount = 12.3
	p.Remittance = ExampleRemittance

	o, err := p.ToQRBytes()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(o))
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
	require.NoError(t, err)

	result, err := p.ToQRBytes()
	require.NoError(t, err)

	expected, err := os.ReadFile("tests/test1.qr")

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	result, err = p.ToQRPNG(QRSize)
	require.NoError(t, err)

	expected, err = os.ReadFile("tests/test1.png")

	require.NoError(t, err)
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
	require.NoError(t, err)

	result, err := p.ToQRBytes()
	require.NoError(t, err)

	expected, err := os.ReadFile("tests/test2.qr")

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	result, err = p.ToQRPNG(QRSize)
	require.NoError(t, err)

	expected, err = os.ReadFile("tests/test2.png")

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
