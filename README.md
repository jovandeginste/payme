# QR code generator (ASCII & PNG) for SEPA payments

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jovandeginste/payme)](https://pkg.go.dev/github.com/jovandeginste/payme)
[![Go Report Card](https://goreportcard.com/badge/github.com/jovandeginste/payme)](https://goreportcard.com/report/github.com/jovandeginste/payme)
[![Go Coverage](https://github.com/jovandeginste/payme/wiki/coverage.svg)](https://raw.githack.com/wiki/jovandeginste/payme/coverage.html)

## What is this?

This tool makes it easier for people (in Europe) to pay with bank transfers. It does not serve as a payment gateway or
anything similar. It does no payment verification at all, it simply provides the necessary payment information in a QR
code, which may be scanned by the person who wants to pay.

The QR code contains the information necessary for a bank transaction in the form of a [SEPA credit
transfer](https://epc-qr.eu/). It can be used to prefill the transaction form if your [banking app supports](#Support)
payment by QR code.

The process of generating the QR code is entirely local and offline. It can be printed in ASCII in the terminal, or
exported as a PNG for inclusion in eg. a web page or a mail.

One QR code can be used without limit, but will always contain the same payment information: amount, remittance message,
destination account. More than one person can scan the same code to pay the same amount (eg. split a bill with friends),
or one person can scan the code on a recurring base (eg. you pay your internet invoice every month and it's a fixed
price)

## Install

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
      --qr-version int         QR code version (default 2)
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

## Support

Please provide feedback if your banking app supports or does not support these QR codes.

The QR code is tested with the mobile apps of these banks:

### Belgium

| Bank    | Support                                    |
| ------- | ------------------------------------------ |
| Belfius | Yes                                        |
| Fortis  | Problems with characters in the remittance |
| KBC     | Yes                                        |

## References

- https://en.wikipedia.org/wiki/Single_Euro_Payments_Area
- https://en.wikipedia.org/wiki/EPC_QR_code
- https://epc-qr.eu/
