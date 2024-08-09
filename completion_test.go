package main

import (
	"bytes"
	"testing"

	"github.com/jovandeginste/payme/payment"
	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	q := qrParams{
		Payment: payment.New(),
	}

	cmdRoot, err := newCommand(&q)
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
	q := qrParams{
		Payment: payment.New(),
	}

	cmdRoot, err := newCommand(&q)
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
	q := qrParams{
		Payment: payment.New(),
	}

	cmdRoot, err := newCommand(&q)
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

func TestVersionCommand(t *testing.T) {
	q := qrParams{
		Payment: payment.New(),
	}

	cmdRoot, err := newCommand(&q)
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
