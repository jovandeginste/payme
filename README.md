Simple CLI tool to generate payment QR codes

Install:

```bash
$ go get github.com/jovandeginste/payme
```

Usage:

```bash
$ payme --help
Generate SEPA payment QR code

Usage:
  payme [flags]

Flags:
      --amount float        Amount of the transaction
      --file string         write code to file, leave empty for stdout
  -h, --help                help for payme
      --iban string         IBAN of the beneficiary
      --name string         Name of the beneficiary
      --output string       output type: png or stdout (default "stdout")
      --remittance string   Remittance (message)
      --structured          Make the remittance (message) structured
```

```bash
$ payme --name "Franz Mustermänn" --iban "DE71110220330123456789" --amount 12.3 --remittance "RF18539007547034"
```

or

```bash
$ payme --name "Franz Mustermänn" --iban "DE71110220330123456789" --amount 12.3 --remittance "RF18539007547034" --output png --file QR.png
```
