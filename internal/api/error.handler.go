package api

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/emarifer/url-shortener-echo-templ-htmx/views/errors_pages"
	"github.com/labstack/echo/v4"
)

func (a *API) CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	var errorPage func(fp bool) templ.Component

	switch code {
	case 401:
		errorPage = errors_pages.Error401
	case 404:
		errorPage = errors_pages.Error404
	case 500:
		errorPage = errors_pages.Error500
	}

	isError = true

	a.renderView(c, errors_pages.ErrorIndex(
		fmt.Sprintf("| Error (%d)", code),
		"",
		fromProtected,
		isError,
		errorPage(fromProtected),
	))
}
