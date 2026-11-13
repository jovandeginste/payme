package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelpCommand(t *testing.T) {
	cmdRoot, err := NewCommand()
	assert.NoError(t, err)

	actualOut := new(bytes.Buffer)
	actualErr := new(bytes.Buffer)

	cmdRoot.SetOut(actualOut)
	cmdRoot.SetErr(actualErr)
	cmdRoot.SetArgs([]string{"help"})

	_, err = cmdRoot.ExecuteC()

	assert.NoError(t, err)
	assert.Contains(t, actualOut.String(), "Generate SEPA payment QR code")
	assert.Empty(t, actualErr.String())
}

func TestCompletionHelp(t *testing.T) {
	cmdRoot, err := NewCommand()
	assert.NoError(t, err)

	actualOut := new(bytes.Buffer)
	actualErr := new(bytes.Buffer)

	cmdRoot.SetOut(actualOut)
	cmdRoot.SetErr(actualErr)
	cmdRoot.SetArgs([]string{"completion", "--help"})

	_, err = cmdRoot.ExecuteC()

	assert.NoError(t, err)
	assert.Contains(t, actualOut.String(), "payme completion [bash|zsh|fish|powershell]")
	assert.Empty(t, actualErr.String())
}

func TestCompletionShells(t *testing.T) {
	cmdRoot, err := NewCommand()
	assert.NoError(t, err)

	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		actualOut := new(bytes.Buffer)
		actualErr := new(bytes.Buffer)

		cmdRoot.SetOut(actualOut)
		cmdRoot.SetErr(actualErr)
		cmdRoot.SetArgs([]string{"completion", shell})

		_, err = cmdRoot.ExecuteC()

		assert.NoError(t, err)
		assert.Contains(t, actualOut.String(), shell+" completion for payme")
		assert.Empty(t, actualErr.String())
	}
}

func TestSixelFileOutput(t *testing.T) {
	file := filepath.Join(t.TempDir(), "qr.six")

	cmdRoot, err := NewCommand()
	require.NoError(t, err)

	cmdRoot.SetArgs([]string{
		"--name", "Franz Mustermänn",
		"--iban", "DE71110220330123456789",
		"--amount", "12.3",
		"--remittance", "RF18539007547034",
		"--output", "sixel",
		"--file", file,
	})

	_, err = cmdRoot.ExecuteC()
	require.NoError(t, err)

	data, err := os.ReadFile(file)
	require.NoError(t, err)

	assert.True(t, bytes.HasPrefix(data, []byte("\x1bPq")))
	assert.True(t, bytes.HasSuffix(data, []byte("\x1b\\")))
}

func TestVersionCommand(t *testing.T) {
	cmdRoot, err := NewCommand()
	assert.NoError(t, err)

	actualOut := new(bytes.Buffer)
	actualErr := new(bytes.Buffer)

	cmdRoot.SetOut(actualOut)
	cmdRoot.SetErr(actualErr)
	cmdRoot.SetArgs([]string{"--version"})

	_, err = cmdRoot.ExecuteC()

	assert.NoError(t, err)
	assert.Contains(t, actualOut.String(), "payme version local (local), built manually")
	assert.Empty(t, actualErr.String())
}
