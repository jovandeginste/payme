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

package kittyimg

import (
	"compress/zlib"
	"encoding/base64"
	"io"
)

const (
	// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#remote-client
	chunkEncSize = 4096
	chunkRawSize = (chunkEncSize / 4) * 3
)

// payloadWriter is an [io.WriteCloser] that encodes the payload binary data in the stream.
// It handles encoding to base64 and 4096 characters chunking.
// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#remote-client
type payloadWriter struct {
	bufEnc [chunkEncSize]byte
	bufRaw [chunkRawSize]byte
	n      int
	w      io.Writer
	ntFrst bool // not first:
}

func (pw *payloadWriter) Reset(w io.Writer) {
	pw.n = 0
	pw.ntFrst = false
	pw.w = w
}

func (pw *payloadWriter) encode() error {
	// fmt.Fprintln(os.Stderr, len(bufRaw), "=>", (len(bufRaw)+2)/3*4)

	base64.StdEncoding.Encode(pw.bufEnc[:], pw.bufRaw[:pw.n])
	_, err := pw.w.Write(pw.bufEnc[:(pw.n+2)/3*4])
	pw.n = 0
	return err
}

func (pw *payloadWriter) Write(b []byte) (n int, err error) {
	for len(b) > 0 {
		if pw.n == cap(pw.bufRaw) {
			if !pw.ntFrst {
				pw.w.Write([]byte{','})
				pw.ntFrst = true
			}
			_, err = pw.w.Write([]byte("m=1;"))
			if err != nil {
				return
			}
			err = pw.encode()
			if err != nil {
				return
			}
			_, err = pw.w.Write([]byte("\033\\\033_G"))
			if err != nil {
				return
			}
		}

		l := copy(pw.bufRaw[pw.n:], b)
		pw.n += l
		n += l
		b = b[l:]
	}
	return
}

// Close closes the writer, flushing any unwritten data to the underlying [io.Writer], but does not close the underlying [io.Writer].
func (pw *payloadWriter) Close() (err error) {
	if pw.n == 0 {
		_, err = pw.w.Write([]byte(";\033\\"))
		return
	}
	_, err = pw.w.Write([]byte{';'})
	if err != nil {
		return
	}
	err = pw.encode()
	if err != nil {
		return
	}
	_, err = pw.w.Write([]byte("\033\\"))
	return
}

// zlibPayloadWriter is an [io.WriteCloser] that adds a [compress/zlib] layer over [payloadWriter].
// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#compression
type zlibPayloadWriter struct {
	buffer [16384]byte
	n      int
	pw     payloadWriter
	zw     *zlib.Writer
}

func (zpw *zlibPayloadWriter) Reset(w io.Writer) {
	_, _ = w.Write([]byte(",o=z"))
	zpw.pw.Reset(w)
	if zpw.zw == nil {
		zpw.zw = zlib.NewWriter(&zpw.pw)
	} else {
		zpw.zw.Reset(&zpw.pw)
	}
	zpw.n = 0
}

func (zpw *zlibPayloadWriter) Write(b []byte) (n int, err error) {
	for len(b) > 0 {
		if zpw.n == cap(zpw.buffer) {
			_, err = zpw.zw.Write(zpw.buffer[:])
			if err != nil {
				return
			}
			zpw.n = 0
		}
		m := copy(zpw.buffer[zpw.n:], b)
		zpw.n += m
		n += m
		b = b[m:]
	}
	return
}

// Close closes the Writer, flushing any unwritten data to the underlying [io.Writer], but does not close the underlying [io.Writer].
func (zp *zlibPayloadWriter) Close() error {
	if zp.n > 0 {
		if _, err := zp.zw.Write(zp.buffer[:zp.n]); err != nil {
			return err
		}
	}
	if err := zp.zw.Close(); err != nil {
		return err
	}
	return zp.pw.Close()
}
