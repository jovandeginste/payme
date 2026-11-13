package payment

import (
	"bytes"
	"io"
	"strings"

	"rsc.io/qr"
)

const halfBlockQuietZone = 4

// ToQRBytes returns an ASCII representation of the QR code
// You can print this to the console, save to a file, etc.
func (p *Payment) ToQRBytes() ([]byte, error) {
	t, err := p.ToString()
	if err != nil {
		return nil, err
	}

	code, err := qr.Encode(t, qr.M)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer

	if err := writeHalfBlocks(&result, code); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func writeHalfBlocks(w io.Writer, code *qr.Code) error {
	const (
		bothDark    = " "
		upperDark   = "▀"
		lowerDark   = "▄"
		bothLight   = "█"
		quietZone   = halfBlockQuietZone
		borderLines = (quietZone + 1) / 2
	)

	border := strings.Repeat(bothLight, code.Size+2*quietZone) + "\n"

	var result strings.Builder

	result.WriteString(strings.Repeat(border, borderLines))

	for y := 0; y < code.Size; y += 2 {
		result.WriteString(strings.Repeat(bothLight, quietZone))

		for x := range code.Size {
			currentDark := code.Black(x, y)
			nextDark := y+1 < code.Size && code.Black(x, y+1)

			switch {
			case currentDark && nextDark:
				result.WriteString(bothDark)
			case currentDark && !nextDark:
				result.WriteString(lowerDark)
			case !currentDark && nextDark:
				result.WriteString(upperDark)
			default:
				result.WriteString(bothLight)
			}
		}

		result.WriteString(strings.Repeat(bothLight, quietZone) + "\n")
	}

	result.WriteString(strings.Repeat(border, borderLines))

	_, err := io.WriteString(w, result.String())

	return err
}
