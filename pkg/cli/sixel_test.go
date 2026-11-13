package cli

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadWithTimeout(t *testing.T) {
	reader, writer := io.Pipe()

	go func() {
		_, _ = writer.Write([]byte("ok"))
	}()

	buf := make([]byte, 2)

	n, err := readWithTimeout(reader, buf, time.Second)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, "ok", string(buf))

	_ = reader.Close()
	_ = writer.Close()
}

func TestReadWithTimeoutExpires(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	defer func() { _ = writer.Close() }()

	buf := make([]byte, 1)

	_, err := readWithTimeout(reader, buf, 10*time.Millisecond)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
}

func TestHasSixelAttribute(t *testing.T) {
	assert.True(t, hasSixelAttribute([]byte("\x1b[?62;4;6;22c")))
	assert.True(t, hasSixelAttribute([]byte("\x1b[?4c")))
	assert.True(t, hasSixelAttribute([]byte("noise\x1b[?62;4;6cmore")))
	assert.False(t, hasSixelAttribute([]byte("\x1b[?62;6;22c")))
	assert.False(t, hasSixelAttribute([]byte("\x1b[?62;22c")))
	assert.False(t, hasSixelAttribute([]byte("garbage")))
	assert.False(t, hasSixelAttribute(nil))
}
