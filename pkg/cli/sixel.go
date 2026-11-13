package cli

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/term"
)

const terminalQueryTimeout = 500 * time.Millisecond

// isSixelSupported asks the terminal for its primary device attributes and
// checks whether it advertises sixel support.
func isSixelSupported() bool {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return false
	}

	defer func() {
		if err := term.Restore(fd, state); err != nil {
			slog.Debug("failed to restore terminal", "error", err)
		}
	}()

	if _, err := os.Stdout.WriteString("\x1b[c"); err != nil {
		return false
	}

	buf := make([]byte, 128)

	n, err := readWithTimeout(os.Stdin, buf, terminalQueryTimeout)
	if err != nil {
		return false
	}

	return hasSixelAttribute(buf[:n])
}

func readWithTimeout(r io.Reader, buf []byte, timeout time.Duration) (int, error) {
	type result struct {
		n   int
		err error
	}

	resultCh := make(chan result, 1)

	go func() {
		n, err := r.Read(buf)
		resultCh <- result{n: n, err: err}
	}()

	select {
	case res := <-resultCh:
		return res.n, res.err
	case <-time.After(timeout):
		return 0, os.ErrDeadlineExceeded
	}
}

func hasSixelAttribute(response []byte) bool {
	for {
		start := bytes.Index(response, []byte("\x1b[?"))
		if start < 0 {
			return false
		}

		response = response[start+3:]

		end := bytes.IndexByte(response, 'c')
		if end < 0 {
			return false
		}

		if slices.Contains(strings.Split(string(response[:end]), ";"), "4") {
			return true
		}

		response = response[end+1:]
	}
}
