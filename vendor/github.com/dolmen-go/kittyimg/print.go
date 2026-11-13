/*
   Copyright 2021-2026 Olivier Mengué.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package kittyimg provides utilities to show image in a graphic terminal emulator supporting [kitty's "terminal graphics protocol"].
//
// [kitty's "terminal graphics protocol"]: https://sw.kovidgoyal.net/kitty/graphics-protocol.html
package kittyimg

import (
	"bytes"
	"fmt"
	"image"
	"io"
)

// Encoder is an [image.Image] encoder, like [image/png.Encoder].
//
// Using an Encoder explicitely allows to reuse resources (memory buffers) when encoding
// multiple images sequentially.
type Encoder struct {
	pw  zlibPayloadWriter
	buf []byte
}

// Encode [encodes] img and writes the result on w.
//
// [encodes]: https://sw.kovidgoyal.net/kitty/graphics-protocol/#display-images-on-screen
func (enc *Encoder) Encode(w io.Writer, img image.Image) error {
	bounds := img.Bounds()

	// f=32 => RGBA
	_, err := fmt.Fprintf(w, "\033_Gq=1,a=T,f=32,s=%d,v=%d,t=d", bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	enc.pw.Reset(w)

	bufCap := min(bounds.Dx()*bounds.Dy()*4, 16384) // Multiple of 4 (RGBA)
	buf := enc.buf
	if cap(enc.buf) < bufCap {
		buf = make([]byte, 0, bufCap)
		enc.buf = buf
	} else {
		buf = buf[:0]
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if len(buf) == cap(buf) {
				if _, err = enc.pw.Write(buf); err != nil {
					return err
				}
				buf = buf[:0]
			}
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			buf = append(buf, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}

	if _, err = enc.pw.Write(buf); err != nil {
		return err
	}
	return enc.pw.Close()
}

// Fprint [encodes] img and writes the result on w.
//
// [encodes]: https://sw.kovidgoyal.net/kitty/graphics-protocol/#display-images-on-screen
func Fprint(w io.Writer, img image.Image) error {
	var e Encoder
	return e.Encode(w, img)
}

// Fprintln calls [Fprint], then writes '\n'.
func Fprintln(w io.Writer, img image.Image) error {
	err := Fprint(w, img)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}

// Transcode transforms the image file into the Kitty protocol representation for display
// on a terminal.
//
// The supported input image formats depend on the formats registered with the [image]
// framework (see [image/png], [image/gif], [image/jpeg]).
func (enc *Encoder) Transcode(w io.Writer, r io.Reader) error {
	var buf bytes.Buffer
	in := io.TeeReader(r, &buf)
	cfg, format, err := image.DecodeConfig(in)
	if err != nil {
		return readError(r, err)
	}
	// Restart from byte 0
	in = io.MultiReader(&buf, r)

	// For PNG we send the raw file that probably has better compression
	// https://sw.kovidgoyal.net/kitty/graphics-protocol/#png-data
	if format == "png" {
		if _, err = fmt.Fprintf(w, "\033_Gq=1,a=T,f=100,s=%d,v=%d", cfg.Width, cfg.Height); err != nil {
			return err
		}

		var pw *payloadWriter = &enc.pw.pw
		pw.Reset(w)

		if _, err = io.Copy(pw, in); err != nil {
			return err
		}
		return pw.Close()
	}

	img, _, err := image.Decode(in)
	if err != nil {
		return readError(r, err)
	}
	return enc.Encode(w, img)
}

// Transcode transforms the image file into the Kitty protocol representation for display
// on a terminal.
//
// The supported input image formats depend on the formats registered with the [image]
// framework (see [image/png], [image/gif], [image/jpeg]).
func Transcode(w io.Writer, r io.Reader) error {
	var enc Encoder
	return enc.Transcode(w, r)
}

func readError(r io.Reader, err error) error {
	if r, ok := r.(interface{ Name() string }); ok {
		if name := r.Name(); name != "" {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return err
}
