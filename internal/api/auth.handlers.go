package api

import (
	"fmt"
	"net/http"

	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/service"
	"github.com/emarifer/url-shortener-echo-templ-htmx/views/auth_views"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

/********** Handler for Home View **********/

func (a *API) homeHandler(c echo.Context) error {
	homeView := auth_views.Home(fromProtected)
	isError = false

	return a.renderView(c, auth_views.HomeIndex(
		"| Home",
		username,
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		homeView,
	))
}

/********** Handlers for Auth Views **********/

func (a *API) registerHandler(c echo.Context) error {
	registerView := auth_views.Register(fromProtected)
	isError = false

	ctx := c.Request().Context()

	if c.Request().Method == "POST" {
		err := a.serv.RegisterUser(
			ctx,
			c.FormValue("email"),
			c.FormValue("username"),
			c.FormValue("password"),
		)
		if err != nil {
			if err == service.ErrUserAlreadyExists {
				setFlashmessages(
					c, "error", service.ErrUserAlreadyExists.Error(),
				)

				return c.Redirect(http.StatusSeeOther, "/register")
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "You have successfully registered!!")

		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return a.renderView(c, auth_views.RegisterIndex(
		"| Register",
		username,
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		registerView,
	))
}

func (a *API) loginHandler(c echo.Context) error {
	loginView := auth_views.Login(fromProtected)
	isError = false

	ctx := c.Request().Context()

	if c.Request().Method == "POST" {
		// obtaining the time zone from the POST request of the login form
		tzone := ""
		if len(c.Request().Header["X-Timezone"]) != 0 {
			tzone = c.Request().Header["X-Timezone"][0]
		}

		// Authentication goes here
		user, err := a.serv.LoginUser(
			ctx,
			c.FormValue("email"),
			c.FormValue("password"),
		)
		if err != nil {
			if err == service.ErrInvalidCredentials {
				setFlashmessages(
					c, "error", service.ErrInvalidCredentials.Error(),
				)

				return c.Redirect(http.StatusSeeOther, "/login")
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		// Get Session and setting Cookies
		sess, _ := session.Get(auth_sessions_key, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600, // in seconds = 1 hour
			HttpOnly: true,
		}

		// Set user as authenticated, their username,
		// their ID and the client's time zone
		sess.Values = map[interface{}]interface{}{
			auth_key:     true,
			user_id_key:  user.UserID,
			username_key: user.Username,
			tzone_key:    tzone,
		}
		sess.Save(c.Request(), c.Response())

		setFlashmessages(c, "success", "You have successfully logged in!!")

		return c.Redirect(http.StatusSeeOther, "/dash")
	}

	return a.renderView(c, auth_views.LoginIndex(
		"| Login",
		username,
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		loginView,
	))
}
