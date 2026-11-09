package main

import (
	"log/slog"
	"os"

	"github.com/jovandeginste/payme/pkg/cli"
)

func main() {
	cmd, err := cli.NewCommand()
	if err != nil {
		slog.Error("failed to create command", "error", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		slog.Error("failed to execute command", "error", err)
		os.Exit(1)
	}
}
