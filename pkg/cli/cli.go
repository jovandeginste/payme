package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jovandeginste/payme/pkg/payment"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI to generate SEPA payment QR codes, either as ASCII or PNG

const qrSize = 300

var (
	gitRefName = "local"
	gitCommit  = "local"
	buildTime  = "manually"
)

type qrParams struct {
	Payment    *payment.Payment
	OutputType string
	OutputFile string
	Debug      bool
}

// NewCommand returns the root command for the payme CLI.
func NewCommand() (*cobra.Command, error) {
	return newCommand(newQRParams())
}

func newQRParams() *qrParams {
	//nolint:exhaustruct_v5 // optional fields are initialized from flag defaults
	return &qrParams{Payment: payment.New()}
}

func newCommand(q *qrParams) (*cobra.Command, error) {
	//nolint:exhaustruct_v5 // cobra defines the remaining fields
	cmdRoot := &cobra.Command{
		Use:           "payme",
		Version:       fmt.Sprintf("%s (%s), built %s\n", gitRefName, gitCommit, buildTime),
		Short:         "Generate SEPA payment QR code",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return q.generate()
		},
	}

	cmdRoot.AddCommand(completionCmd(cmdRoot))

	if err := q.init(cmdRoot); err != nil {
		return nil, err
	}

	return cmdRoot, nil
}

func (q *qrParams) init(cmdRoot *cobra.Command) error {
	viper.SetEnvPrefix("PAYME")

	for _, e := range []string{"name", "bic", "iban"} {
		if err := viper.BindEnv(e); err != nil {
			return err
		}
	}

	cmdRoot.Flags().StringVar(&q.OutputType, "output", "stdout", "output type: stdout, png, sixel or kitty")
	cmdRoot.Flags().StringVar(&q.OutputFile, "file", "", "write code to file, leave empty for stdout")
	cmdRoot.Flags().BoolVar(&q.Debug, "debug", false, "print debug output")

	cmdRoot.Flags().IntVar(&q.Payment.CharacterSet, "character-set", 2, "QR code character set")
	cmdRoot.Flags().IntVar(&q.Payment.Version, "qr-version", 2, "QR code version")
	cmdRoot.Flags().StringVar(&q.Payment.NameBeneficiary, "name", viper.GetString("name"), "Name of the beneficiary")
	cmdRoot.Flags().StringVar(&q.Payment.BICBeneficiary, "bic", viper.GetString("bic"), "BIC of the beneficiary")
	cmdRoot.Flags().StringVar(&q.Payment.IBANBeneficiary, "iban", viper.GetString("iban"), "IBAN of the beneficiary")
	cmdRoot.Flags().Float64Var(&q.Payment.EuroAmount, "amount", 0, "Amount of the transaction")
	cmdRoot.Flags().StringVar(&q.Payment.Remittance, "remittance", "", "Remittance (message)")
	cmdRoot.Flags().StringVar(&q.Payment.Purpose, "purpose", "", "Purpose of the transaction")
	cmdRoot.Flags().BoolVar(&q.Payment.RemittanceIsStructured, "structured", false, "Make the remittance (message) structured")

	return nil
}

func (q *qrParams) generate() error {
	if q.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("parameters", "params", fmt.Sprintf("%#v", q))
	}

	var (
		qr  []byte
		err error
	)

	switch q.OutputType {
	case "kitty":
		qr, err = q.generateQRKitty()
	case "png":
		qr, err = q.generateQRPNG()
	case "sixel":
		qr, err = q.generateQRSixel()
	case "stdout":
		qr, err = q.generateQRStdout()
	default:
		return fmt.Errorf("unknown output type: %q", q.OutputType)
	}

	if err != nil {
		return err
	}

	if q.OutputFile == "" {
		_, err = fmt.Fprint(os.Stdout, string(qr))
		return err
	}

	return os.WriteFile(q.OutputFile, qr, 0o600)
}

func (q *qrParams) generateQRStdout() ([]byte, error) {
	if q.Debug {
		s, err := q.Payment.ToString()
		if err != nil {
			return nil, err
		}

		slog.Debug("payment data", "data", s)
	}

	return q.Payment.ToQRBytes()
}

func (q *qrParams) generateQRSixel() ([]byte, error) {
	if q.Debug {
		s, err := q.Payment.ToString()
		if err != nil {
			return nil, err
		}

		slog.Debug("payment data", "data", s)
	}

	if q.OutputFile == "" && !isSixelSupported() {
		return nil, fmt.Errorf("terminal does not support sixel graphics; use --output stdout or --file")
	}

	return q.Payment.ToQRSixel()
}

func (q *qrParams) generateQRKitty() ([]byte, error) {
	if q.Debug {
		s, err := q.Payment.ToString()
		if err != nil {
			return nil, err
		}

		slog.Debug("payment data", "data", s)
	}

	return q.Payment.ToQRKitty(qrSize)
}

func (q *qrParams) generateQRPNG() ([]byte, error) {
	if q.Debug {
		s, err := q.Payment.ToString()
		if err != nil {
			return nil, err
		}

		slog.Debug("payment data", "data", s)
	}

	return q.Payment.ToQRPNG(qrSize)
}
