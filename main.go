package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jovandeginste/payme/payment"
	"github.com/spf13/cobra"
)

const QRSize = 300

type QRParams struct {
	Payment    *payment.Payment
	OutputType string
	OutputFile string
	Debug      bool
}

func main() {
	q := QRParams{
		Payment: payment.New(),
	}

	cmdRoot := &cobra.Command{
		Use:   "payme",
		Short: "Generate SEPA payment QR code",
		Run: func(cmd *cobra.Command, args []string) {
			q.generate()
		},
	}

	cmdRoot.Flags().StringVar(&q.OutputType, "output", "stdout", "output type: png or stdout")
	cmdRoot.Flags().StringVar(&q.OutputFile, "file", "", "write code to file, leave empty for stdout")
	cmdRoot.Flags().BoolVar(&q.Debug, "debug", false, "print debug output")

	cmdRoot.Flags().IntVar(&q.Payment.CharacterSet, "character-set", 2, "QR code character set")
	cmdRoot.Flags().StringVar(&q.Payment.NameBeneficiary, "name", "", "Name of the beneficiary")
	cmdRoot.Flags().StringVar(&q.Payment.BICBeneficiary, "bic", "", "BIC of the beneficiary")
	cmdRoot.Flags().IntVar(&q.Payment.Version, "version", 2, "QR code version")
	cmdRoot.Flags().StringVar(&q.Payment.IBANBeneficiary, "iban", "", "IBAN of the beneficiary")
	cmdRoot.Flags().Float64Var(&q.Payment.EuroAmount, "amount", 0, "Amount of the transaction")
	cmdRoot.Flags().StringVar(&q.Payment.Remittance, "remittance", "", "Remittance (message)")
	cmdRoot.Flags().StringVar(&q.Payment.Purpose, "purpose", "", "Purpose of the transaction")
	cmdRoot.Flags().BoolVar(&q.Payment.RemittanceIsStructured, "structured", false, "Make the remittance (message) structured")

	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}

func (q *QRParams) generate() {
	var (
		qr  []byte
		err error
	)

	if q.Debug {
		log.Printf("%#v\n", q)
	}

	switch q.OutputType {
	case "png":
		qr, err = q.generateQRPNG()
	case "stdout":
		qr, err = q.generateQRStdout()
	}

	if err != nil {
		log.Fatal(err)
	}

	if q.OutputFile == "" {
		fmt.Fprintf(os.Stdout, "%s", qr)
		return
	}

	err = ioutil.WriteFile(q.OutputFile, qr, 0o600)
	if err != nil {
		log.Fatal(err)
	}
}

func (q *QRParams) generateQRStdout() ([]byte, error) {
	p := q.Payment

	if q.Debug {
		s, err := p.ToString()
		if err != nil {
			return nil, err
		}

		log.Print("Data: ", s)
	}

	return p.ToQRString()
}

func (q *QRParams) generateQRPNG() ([]byte, error) {
	p := q.Payment

	if q.Debug {
		s, err := p.ToString()
		if err != nil {
			return nil, err
		}

		log.Print("Data: ", s)
	}

	return p.ToQRPNG(QRSize)
}
