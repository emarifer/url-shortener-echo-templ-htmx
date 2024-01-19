package api

import "github.com/labstack/echo/v4"

func (a *API) RegisterRoutes(e *echo.Echo) {
	e.GET("/", a.homeHandler)
	e.GET("/register", a.registerHandler)
	e.POST("/register", a.registerHandler)
	e.GET("/login", a.loginHandler)
	e.POST("/login", a.loginHandler)
	e.GET("/s/:slug", a.linkRedirectHandler)

	protectedGroup := e.Group("/dash", a.authMiddleware)
	/* ↓ Protected Routes ↓ */
	protectedGroup.GET("", a.dashboardHandler)
	protectedGroup.GET("/create", a.createLinkHandler)
	protectedGroup.POST("/create", a.createLinkHandler)
	protectedGroup.GET("/randomize", a.createRandomSlugHandler)
	protectedGroup.GET("/edit", a.updateLinkHandler)
	protectedGroup.POST("/edit", a.updateLinkHandler)
	protectedGroup.DELETE("/delete/:slug", a.deleteLinkHandler)
	protectedGroup.POST("/logout", a.logoutHandler)
}
