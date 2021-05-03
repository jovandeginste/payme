# echo-logrusmiddleware

![logrus middleware](/logrus.png)

An adapter (middleware) to make the Golang [Echo web
framework](https://github.com/labstack/echo) logging work with
[logrus](https://github.com/sirupsen/logrus), an excellent logging solution.

Improves upon [sandalwing/echo-logrusmiddleware](https://github.com/sandalwing/echo-logrusmiddleware) by:
1. Using the correct import for logrus
2. Including the request_id prop in the log output in order to support Echo's request ID middleware.
3. Supporting Echo v4

## Install

```
$ go get github.com/alexferl/echo-logrusmiddleware
```

## Usage

```go
package main

import (
	"github.com/alexferl/echo-logrusmiddleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	// echo Logger interface friendly wrapper around logrus logger to use it
	// for default echo logger
	e.Logger = logrusmiddleware.Logger{logrus.StandardLogger()}
	e.Use(logrusmiddleware.Hook())

	// do the rest of your echo setup, routes, listen on server, etc..
}
```
