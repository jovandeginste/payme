package main

import (
	"flag"
	"fmt"
	"log"
)

const QRSize = 300

type QRParams struct {
	NameBeneficiary string
	IBANBeneficiary string
	Amount          float64
	Remittance      string
	IsStructured    bool
	OutputType      string
}

func main() {
	p := new(QRParams)

	flag.StringVar(&p.NameBeneficiary, "name", "", "Name of the beneficiary")
	flag.StringVar(&p.IBANBeneficiary, "iban", "", "IBAN of the beneficiary")
	flag.Float64Var(&p.Amount, "amount", 0, "Amount of the transaction")
	flag.StringVar(&p.Remittance, "remittance", "", "Remittance (message)")
	flag.BoolVar(&p.IsStructured, "structured", false, "Make the remittance (message) structured")
	flag.StringVar(&p.OutputType, "output", "stdout", "output type: png or stdout")

	flag.Parse()

	log.Printf("%#v\n", p)

	switch p.OutputType {
	case "png":
		qr, err := generateQRPNG(p)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(qr)
	case "stdout":
		err := generateQRStdout(p)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateQRStdout(params *QRParams) error {
	p := NewPayment()

	p.NameBeneficiary = params.NameBeneficiary
	p.IBANBeneficiary = params.IBANBeneficiary
	p.EuroAmount = params.Amount
	p.Remittance = params.Remittance
	p.RemittanceIsStructured = params.IsStructured

	return p.ToQRStdout()
}

func generateQRPNG(params *QRParams) (string, error) {
	p := NewPayment()

	p.NameBeneficiary = params.NameBeneficiary
	p.IBANBeneficiary = params.IBANBeneficiary
	p.EuroAmount = params.Amount
	p.Remittance = params.Remittance
	p.RemittanceIsStructured = params.IsStructured

	png, err := p.ToQRPNG(QRSize)
	if err != nil {
		return "", err
	}

	return string(png), nil
}
