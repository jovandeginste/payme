Simple CLI tool to generate payment QR codes

Install:

```bash
$ go get github.com/jovandeginste/payme
```

Usage:

```bash
$ payme --help
  -amount float
    	Amount of the transaction
  -iban string
    	IBAN of the beneficiary
  -name string
    	Name of the beneficiary
  -remittance string
    	Remittance (message)
  -structured
    	Make the remittance (message) structured
```
