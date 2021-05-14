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

var p QRParams

type QRParams struct {
	NameBeneficiary string
	IBANBeneficiary string
	Amount          float64
	Remittance      string
	IsStructured    bool
	OutputType      string
	OutputFile      string
	Debug           bool
}

func main() {
	cmdRoot := &cobra.Command{
		Use:   "payme",
		Short: "Generate SEPA payment QR code",
		Run: func(cmd *cobra.Command, args []string) {
			generate()
		},
	}

	cmdRoot.Flags().StringVar(&p.NameBeneficiary, "name", "", "Name of the beneficiary")
	cmdRoot.Flags().StringVar(&p.IBANBeneficiary, "iban", "", "IBAN of the beneficiary")
	cmdRoot.Flags().Float64Var(&p.Amount, "amount", 0, "Amount of the transaction")
	cmdRoot.Flags().StringVar(&p.Remittance, "remittance", "", "Remittance (message)")
	cmdRoot.Flags().BoolVar(&p.IsStructured, "structured", false, "Make the remittance (message) structured")
	cmdRoot.Flags().StringVar(&p.OutputType, "output", "stdout", "output type: png or stdout")
	cmdRoot.Flags().StringVar(&p.OutputFile, "file", "", "write code to file, leave empty for stdout")
	cmdRoot.Flags().BoolVar(&p.Debug, "debug", false, "print debug output")

	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}

func generate() {
	var (
		qr  []byte
		err error
	)

	if p.Debug {
		log.Printf("%#v\n", p)
	}

	switch p.OutputType {
	case "png":
		qr, err = generateQRPNG(p)
	case "stdout":
		qr, err = generateQRStdout(p)
	}

	if err != nil {
		log.Fatal(err)
	}

	if p.OutputFile == "" {
		fmt.Fprintf(os.Stdout, "%s", qr)
		return
	}

	err = ioutil.WriteFile(p.OutputFile, qr, 0o600)
	if err != nil {
		log.Fatal(err)
	}
}

func (params *QRParams) preparePayment() payment.Payment {
	p := payment.New()

	p.NameBeneficiary = params.NameBeneficiary
	p.IBANBeneficiary = params.IBANBeneficiary
	p.EuroAmount = params.Amount
	p.Remittance = params.Remittance
	p.RemittanceIsStructured = params.IsStructured

	return p
}

func generateQRStdout(params QRParams) ([]byte, error) {
	p := params.preparePayment()

	if params.Debug {
		s, err := p.ToString()
		if err != nil {
			return nil, err
		}

		log.Print("Data: ", s)
	}

	return p.ToQRString()
}

func generateQRPNG(params QRParams) ([]byte, error) {
	p := params.preparePayment()

	if params.Debug {
		s, err := p.ToString()
		if err != nil {
			return nil, err
		}

		log.Print("Data: ", s)
	}

	return p.ToQRPNG(QRSize)
}
