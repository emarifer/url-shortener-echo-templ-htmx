package api

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (a *API) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			// fmt.Println(ok, auth)
			fromProtected = false
			username = ""

			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Please provide valid credentials")
		}

		if userId, ok := sess.Values[user_id_key].(string); ok && len(userId) != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if uNameFromSession, ok := sess.Values[username_key].(string); ok && len(uNameFromSession) != 0 {
			username = uNameFromSession
			// set the username in the context
			c.Set(username_key, uNameFromSession)
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		fromProtected = true

		return next(c)
	}
}
