package api

import (
	"github.com/a-h/templ"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// keys for auth session & Echo context (userId, username and time zone)
// (in production, these values must be in environment variables)
const (
	auth_sessions_key string = "authenticate-sessions"
	auth_key          string = "authenticated"
	user_id_key       string = "user_id"
	username_key      string = "username"
	tzone_key         string = "time_zone"
	secret_key        string = "01234567890123456789012345678901"
)

// flags for views
var (
	fromProtected bool   = false
	isError       bool   = false
	username      string = ""
)

type API struct {
	serv          service.Service
	dataValidator *validator.Validate
}

func New(serv service.Service) *API {

	return &API{
		serv:          serv,
		dataValidator: validator.New(),
	}
}

func (a *API) Start(e *echo.Echo, address string) error {
	e.Static("/", "assets")
	// set error handler
	e.HTTPErrorHandler = a.CustomHTTPErrorHandler
	// session middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(secret_key))))

	a.RegisterRoutes(e)

	return e.Start(address)
}

func (a *API) renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
