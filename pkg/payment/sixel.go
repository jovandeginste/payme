package payment

import (
	"bytes"
	"fmt"
	"io"

	"rsc.io/qr"
)

const (
	sixelBlockSize = 6
	sixelQuietZone = 4
)

// ToQRSixel returns a sixel representation of the QR code
// Use this when the terminal supports sixel graphics.
func (p *Payment) ToQRSixel() ([]byte, error) {
	t, err := p.ToString()
	if err != nil {
		return nil, err
	}

	code, err := qr.Encode(t, qr.M)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer

	if err := writeSixel(&result, code); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func writeSixel(w io.Writer, code *qr.Code) error {
	width := sixelBlockSize * (code.Size + 2*sixelQuietZone)
	border := fmt.Appendf(nil, "#1!%d~-\n", width)

	if _, err := w.Write([]byte("\x1bPq\n#0;2;0;0;0#1;2;100;100;100\n")); err != nil {
		return err
	}

	for range sixelQuietZone {
		if _, err := w.Write(border); err != nil {
			return err
		}
	}

	for y := range code.Size {
		if err := writeSixelRow(w, code, y); err != nil {
			return err
		}
	}

	for range sixelQuietZone {
		if _, err := w.Write(border); err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\x1b\\"))

	return err
}

func writeSixelRow(w io.Writer, code *qr.Code, y int) error {
	row := fmt.Appendf(nil, "#1!%d~", sixelBlockSize*sixelQuietZone)

	last := -1
	repeat := 0

	flush := func() {
		if repeat == 0 {
			return
		}

		row = fmt.Appendf(row, "#%d!%d~", last, sixelBlockSize*repeat)
		repeat = 0
	}

	for x := range code.Size {
		color := 1
		if code.Black(x, y) {
			color = 0
		}

		if color != last {
			flush()

			last = color
		}

		repeat++
	}

	flush()

	row = fmt.Appendf(row, "#1!%d~-\n", sixelBlockSize*sixelQuietZone)

	_, err := w.Write(row)

	return err
}
