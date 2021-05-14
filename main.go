package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jovandeginste/payme/payment"
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

	var (
		qr  string
		err error
	)

	switch p.OutputType {
	case "png":
		qr, err = generateQRPNG(p)
	case "stdout":
		qr, err = generateQRStdout(p)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(qr)
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

func generateQRStdout(params *QRParams) (string, error) {
	p := params.preparePayment()

	return p.ToQRString()
}

func generateQRPNG(params *QRParams) (string, error) {
	p := params.preparePayment()

	png, err := p.ToQRPNG(QRSize)
	if err != nil {
		return "", err
	}

	return string(png), nil
}
