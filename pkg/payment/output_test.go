package payment_test

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/jovandeginste/payme/pkg/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	QRSize            = 300
	ExampleIBAN       = "FR1420041010050500013M02606"
	ExampleIBANDE     = "DE71110220330123456789"
	ExampleName       = "François D'Alsace S.A."
	ExampleRemittance = "Client:Marie Louise La Lune"
)

// assertPNGEqual compares the decoded content of two PNG images. The encoded
// bytes are not compared, since they depend on the PNG encoder version.
func assertPNGEqual(t *testing.T, expected, actual []byte) {
	t.Helper()

	expectedImage, err := png.Decode(bytes.NewReader(expected))
	require.NoError(t, err)

	actualImage, err := png.Decode(bytes.NewReader(actual))
	require.NoError(t, err)

	require.Equal(t, expectedImage.Bounds(), actualImage.Bounds())

	bounds := expectedImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, a1 := expectedImage.At(x, y).RGBA()
			r2, g2, b2, a2 := actualImage.At(x, y).RGBA()

			require.Equal(t, [4]uint32{r1, g1, b1, a1}, [4]uint32{r2, g2, b2, a2},
				"pixel (%d, %d) differs", x, y)
		}
	}
}

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
	assertPNGEqual(t, expected, result)
}

func TestStructuredPaymentQR(t *testing.T) {
	p := payment.NewStructured()

	p.Version = 1
	p.CharacterSet = 1
	p.BICBeneficiary = "BHBLDEHHXXX"
	p.NameBeneficiary = "Franz Mustermänn"
	p.IBANBeneficiary = ExampleIBANDE
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
	assertPNGEqual(t, expected, result)
}
