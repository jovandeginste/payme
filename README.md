# QR code generator (ASCII & PNG) for SEPA payments

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jovandeginste/payme)](https://pkg.go.dev/github.com/jovandeginste/payme)
[![Go Report Card](https://goreportcard.com/badge/github.com/jovandeginste/payme)](https://goreportcard.com/report/github.com/jovandeginste/payme)
[![Coverage](http://gocover.io/_badge/github.com/jovandeginste/payme)](http://gocover.io/github.com/jovandeginste/payme)

Install:

```bash
$ go install github.com/jovandeginste/payme@latest
```

Usage:

```bash
Generate SEPA payment QR code

Usage:
  payme [flags]

Flags:
      --amount float        Amount of the transaction
      --bic string          BIC of the beneficiary
      --character-set int   QR code character set (default 2)
      --debug               print debug output
      --file string         write code to file, leave empty for stdout
  -h, --help                help for payme
      --iban string         IBAN of the beneficiary
      --name string         Name of the beneficiary
      --output string       output type: png or stdout (default "stdout")
      --purpose string      Purpose of the transaction
      --remittance string   Remittance (message)
      --structured          Make the remittance (message) structured
      --version int         QR code version (default 2)
```

You can set some default values in your ENV, eg.:

```bash
export PAYME_IBAN=DE71110220330123456789
export PAYME_NAME="Franz Mustermänn"
export PAYME_BIC=BHBLDEHHXXX
```

Generate QR code as text, print on the console:

```bash
$ payme \
  --name "Franz Mustermänn" \
  --iban "DE71110220330123456789" \
  --amount 12.3 \
  --remittance "RF18539007547034"
```

Generate QR code as png, save as file:

```bash
$ payme \
  --name "Franz Mustermänn" \
  --iban "DE71110220330123456789" \
  --amount 12.3 \
  --remittance "RF18539007547034" \
  --output png \
  --file QR.png
```
