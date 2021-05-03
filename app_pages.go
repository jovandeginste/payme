package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	echo "github.com/labstack/echo/v4"
)

type QRParams struct {
	NameBeneficiary string  `query:"name_beneficiary"`
	IBANBeneficiary string  `query:"iban_beneficiary"`
	Amount          float64 `query:"amount"`
	Remittance      string  `query:"remittance"`
	IsStructured    bool    `query:"structured"`
}

//go:embed assets
var embededFiles embed.FS

func (a *App) GetFileSystem() http.FileSystem {
	if a.Live {
		a.Logger.Debug("using live mode")

		return http.FS(os.DirFS("assets"))
	}

	a.Logger.Debug("using embed mode")

	fsys, err := fs.Sub(embededFiles, "assets")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func (a *App) GeneratePNG(c echo.Context) error {
	p := new(QRParams)
	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	a.Logger.Infof("%#v", p)

	qr, err := a.GenerateQR(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	a.Logger.Infof("%#v", err)

	return c.String(http.StatusOK, qr)
}

func (a *App) GenerateQR(params *QRParams) (string, error) {
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
