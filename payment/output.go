package payment

import (
	"bytes"
	"image/png"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/mdp/qrterminal/v3"
)

func (p *Payment) ToString() (string, error) {
	if err := p.ValidateFields(); err != nil {
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

func (p *Payment) BICBeneficiaryString() string {
	if p.Version != 1 {
		return ""
	}

	return p.BICBeneficiary
}

func (p *Payment) ToQRString() ([]byte, error) {
	var result bytes.Buffer

	t, err := p.ToString()
	if err != nil {
		return nil, err
	}

	qrterminal.GenerateHalfBlock(t, qrterminal.M, &result)

	return result.Bytes(), nil
}

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
