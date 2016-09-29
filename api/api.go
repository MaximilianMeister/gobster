// Package api provides a http interface to feed
package api

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	"github.com/MaximilianMeister/gobster/feed"
)

// Serve starts the webserver
func Serve() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello!")
	})

	e.GET("/:name", func(c echo.Context) error {
		quote, err := feed.Get(c.Param("name"))
		if err != nil {
			return c.String(http.StatusNotAcceptable, "Not Acceptable")
		}

		return c.String(http.StatusOK, quote)
	})
	e.POST("/:name", func(c echo.Context) error {
		err := feed.Set(c.Param("name"), c.FormValue("quote"))
		if err != nil {
			return c.String(http.StatusNotAcceptable, "Not Acceptable")
		}

		return c.String(http.StatusOK, "Danke")
	})
	e.GET(fmt.Sprintf("/:name/all"), func(c echo.Context) error {
		quotes, err := feed.GetAll(c.Param("name"))
		if err != nil {
			return c.String(http.StatusNotAcceptable, "Not Acceptable")
		}

		var buffer bytes.Buffer
		for i, n := range quotes {
			buffer.WriteString(fmt.Sprintf("%d. %s\n", i, n))
		}

		return c.String(http.StatusOK, buffer.String())
	})

	e.Run(standard.New(":9876"))
}
