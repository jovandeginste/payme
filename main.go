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
}

func main() {
	p := new(QRParams)

	flag.StringVar(&p.NameBeneficiary, "name", "", "Name of the beneficiary")
	flag.StringVar(&p.IBANBeneficiary, "iban", "", "IBAN of the beneficiary")
	flag.Float64Var(&p.Amount, "amount", 0, "Amount of the transaction")
	flag.StringVar(&p.Remittance, "remittance", "", "Remittance (message)")
	flag.BoolVar(&p.IsStructured, "structured", false, "Make the remittance (message) structured")

	flag.Parse()

	log.Printf("%#v\n", p)

	qr, err := generateQR(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(qr)
}

func generateQR(params *QRParams) (string, error) {
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
