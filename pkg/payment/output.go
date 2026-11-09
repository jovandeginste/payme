package payment

import (
	"bytes"
	"image/png"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/mdp/qrterminal/v3"
)

// ToString returns the content of the QR code as string
// Use this to then generate the QR code in the form you need
func (p *Payment) ToString() (string, error) {
	if err := p.IsValid(); err != nil {
		return "", err
	}

	fields := []string{
		p.ServiceTag,
		p.VersionString(),
		p.CharacterSetString(),
		p.IdentificationCode,
		p.BICBeneficiaryString(),
		p.NameBeneficiary,
		p.IBANBeneficiaryString(),
		p.EuroAmountString(),
		p.PurposeString(),
		p.RemittanceStructured(),
		p.RemittanceText(),
		p.B2OInformation,
	}

	return strings.Join(fields, "\n"), nil
}

// ToQRBytes returns an ASCII representation of the QR code
// You can print this to the console, save to a file, etc.
func (p *Payment) ToQRBytes() ([]byte, error) {
	var result bytes.Buffer

	t, err := p.ToString()
	if err != nil {
		return nil, err
	}

	qrterminal.GenerateHalfBlock(t, qrterminal.M, &result)

	return result.Bytes(), nil
}

// ToQRPNG returns an PNG representation of the QR code
// You should save this to a file, or pass it to an image processing library
func (p *Payment) ToQRPNG(qrSize int) ([]byte, error) {
	t, err := p.ToString()
	if err != nil {
		return nil, err
	}

	// Create the barcode
	qrCode, err := qr.Encode(t, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	// Scale the barcode to qrSize x qrSize pixels
	qrCode, err = barcode.Scale(qrCode, qrSize, qrSize)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	// encode the barcode as png
	err = png.Encode(&b, qrCode)

	return b.Bytes(), err
}
