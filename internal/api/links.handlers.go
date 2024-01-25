package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/emarifer/url-shortener-echo-templ-htmx/encryption"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/api/dto"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/entity"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/service"
	"github.com/emarifer/url-shortener-echo-templ-htmx/views/components"
	"github.com/emarifer/url-shortener-echo-templ-htmx/views/links_views"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/********** Handlers for links actions & dashboard render **********/

func (a *API) dashboardHandler(c echo.Context) error {
	isError = false

	params := dto.SearchQuery{}
	binder := &echo.DefaultBinder{}
	binder.BindQueryParams(c, &params)
	// fmt.Println(params.Search)
	searchDescription := strings.Trim(params.Search, " ")
	searchEvent := c.Request().Header.Get("HX-Trigger-Name")

	ctx := c.Request().Context()
	userId := c.Get(user_id_key).(string)

	if searchEvent == "search" {
		// to observe the request indicator
		time.Sleep(1 * time.Second)

		searchedLinks, err := a.serv.SearchLinksByDescription(
			ctx, searchDescription, userId,
		)
		if err != nil {

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return a.renderView(c, components.CardList(
			c.Scheme()+"://"+c.Request().Host,
			c.Get(tzone_key).(string),
			searchedLinks,
		))
	}

	titlePage := fmt.Sprintf(
		"| %s's Links",
		cases.Title(language.English).String(c.Get(username_key).(string)),
	)

	if searchEvent != "search" && searchDescription != "" {
		searchedLinks, err := a.serv.SearchLinksByDescription(
			ctx, searchDescription, userId,
		)
		if err != nil {

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return a.renderView(c, links_views.DashboardIndex(
			titlePage,
			c.Get(username_key).(string),
			fromProtected,
			isError,
			getFlashmessages(c, "error"),
			getFlashmessages(c, "success"),
			links_views.Dashboard(
				titlePage,
				c.Scheme()+"://"+c.Request().Host,
				c.Get(tzone_key).(string),
				searchedLinks,
			),
		))

	}

	links, err := a.serv.RecoverLinks(ctx, userId)
	if err != nil {

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	return a.renderView(c, links_views.DashboardIndex(
		titlePage,
		c.Get(username_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		links_views.Dashboard(
			titlePage,
			c.Scheme()+"://"+c.Request().Host,
			c.Get(tzone_key).(string),
			links,
		),
	))
}

func (a *API) createLinkHandler(c echo.Context) error {
	isError = false

	ctx := c.Request().Context()

	if c.Request().Method == "POST" {
		link := &entity.Link{
			Url:         strings.Trim(c.FormValue("url"), " "),
			Slug:        strings.Trim(c.FormValue("slug"), " "),
			Description: strings.Trim(c.FormValue("description"), " "),
			UserID:      c.Get(user_id_key).(string),
		}

		err := a.serv.AddLink(ctx, link)
		if err != nil {
			if err == service.ErrDuplicateUrl {
				setFlashmessages(
					c, "error", service.ErrDuplicateUrl.Error(),
				)

				return c.Redirect(http.StatusSeeOther, "/dash/create")
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "Short link created successfully!!")

		return c.Redirect(http.StatusSeeOther, "/dash")
	}

	return a.renderView(c, links_views.DashboardIndex(
		"| Create Link",
		c.Get(username_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		links_views.LinkCreator(""),
	))
}

func (a *API) updateLinkHandler(c echo.Context) error {
	isError = false

	ctx := c.Request().Context()
	params := dto.LinkEdit{}

	if c.Request().Method == "GET" {
		if err := c.Bind(&params); err != nil {

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}
	}

	if c.Request().Method == "POST" {
		binder := &echo.DefaultBinder{}
		binder.BindHeaders(c, &params)
		// fmt.Println("X-Id:", params.ID)

		err := a.serv.UpdateLink(
			ctx,
			strings.Trim(c.FormValue("description"), " "),
			params.ID,
		)
		if err != nil {
			if err == service.ErrResourceNotFound {

				return echo.NewHTTPError(
					echo.ErrNotFound.Code,
					fmt.Sprintf(
						"something went wrong: %s",
						err,
					))
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "Short link updated successfully!!")

		return c.Redirect(http.StatusSeeOther, "/dash")
	}

	return a.renderView(c, components.LinkEditor(
		params.Slug, params.Description, params.ID,
	))
}

func (a *API) deleteLinkHandler(c echo.Context) error {
	isError = false

	ctx := c.Request().Context()
	slug := c.Param("slug")
	userId := c.Get(user_id_key).(string)

	err := a.serv.RemoveLink(ctx, slug, userId)
	if err != nil {
		if err == service.ErrResourceNotFound {

			return echo.NewHTTPError(
				echo.ErrNotFound.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	setFlashmessages(c, "success", "Short link successfully removed!!")

	return c.Redirect(http.StatusSeeOther, "/dash")
}

func (a *API) linkRedirectHandler(c echo.Context) error {
	isError = false

	ctx := c.Request().Context()
	slug := c.Param("slug")

	sl, err := a.serv.RecoverLink(ctx, slug)
	if err != nil {
		if err == service.ErrResourceNotFound {

			return echo.NewHTTPError(
				echo.ErrNotFound.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	return c.Redirect(http.StatusSeeOther, sl.Url)
}

func (a *API) createRandomSlugHandler(c echo.Context) error {
	randomSlug, err := encryption.CreateSlug(6)
	if err != nil {

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	return a.renderView(c, components.InputSlug(randomSlug))
}

/********** Handler for logout action **********/

func (a *API) logoutHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	// Revoke users authentication
	sess.Values = map[interface{}]interface{}{
		auth_key:     false,
		user_id_key:  "",
		username_key: "",
		tzone_key:    "",
	}
	sess.Save(c.Request(), c.Response())

	setFlashmessages(c, "success", "You have successfully logged out!!")

	fromProtected = false

	return c.Redirect(http.StatusSeeOther, "/")
}
