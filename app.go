package main

import (
	"fmt"
	"net/http"

	logrusmiddleware "github.com/alexferl/echo-logrusmiddleware"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type App struct {
	Port int
	Live bool

	Logger *logrus.Logger
	Echo   *echo.Echo
}

func (a *App) Initialize() {
	a.Logger.Infof("Initializing web server...")

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Logger = logrusmiddleware.Logger{Logger: a.Logger}

	e.Use(logrusmiddleware.Hook())
	e.Use(middleware.Recover())

	assetHandler := http.FileServer(a.GetFileSystem())

	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	e.GET("/generate.png", a.GeneratePNG)

	a.Echo = e
}

func (a *App) PortString() string {
	return fmt.Sprintf(":%d", a.Port)
}

func (a *App) StartServer() error {
	a.Logger.Infof("Starting web server on %#v", a.PortString())

	return a.Echo.Start(a.PortString())
}
