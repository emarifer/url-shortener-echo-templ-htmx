package main

import (
	"context"
	"fmt"

	"github.com/emarifer/url-shortener-echo-templ-htmx/database"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/api"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/repository"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/service"
	"github.com/emarifer/url-shortener-echo-templ-htmx/settings"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			settings.New,
			database.NewPostgresDB,
			repository.New,
			service.New,
			echo.New,
			api.New,
		),

		fx.Invoke(
			setLifeCycle,
		),
	)

	app.Run()
}

func setLifeCycle(
	lc fx.Lifecycle,
	a *api.API,
	s *settings.Settings,
	e *echo.Echo,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf(":%s", s.Port)
			go func() {
				e.Logger.Fatal(a.Start(e, address))
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {

			return e.Close()
		},
	})
}

/* REFERENCES:
https://htmx.org/essays/rest-copypasta/
A Customized Confirmation UI (HTMX):
https://htmx.org/examples/confirm/

HTML <input> pattern Attribute:
https://www.w3schools.com/tags/att_input_pattern.asp

How to Use Tailwind CSS Grid:
https://refine.dev/blog/tailwind-grid/#introduction
Container (Tailwind Docs):
https://tailwindcss.com/docs/container
How to hide scrollbar on your element in Tailwind:
https://dev.to/derick1530/how-to-create-scrollable-element-in-tailwind-without-a-scrollbar-4mbd
Can't center absolute position (Tailwind.css):
https://stackoverflow.com/questions/60362442/cant-center-absolute-position-tailwind-css

Binding Multiple Sources & Binding Direct Source:
https://echo.labstack.com/docs/binding#multiple-sources
https://stackoverflow.com/questions/72875958/golang-binding-headers-in-echo-api

Postgresql - Composite unique constraints:
https://www.beekeeperstudio.io/blog/guide-to-unique-constraints-in-postgresql#composite-unique-constraints

JavaScript to toggle "text/password" in
the view password button of the login form:
onclick="
	const e = document.getElementById('togglePassword');
	if (e.type == 'password') {
		e.type = 'text';
	} else {
		e.type = 'password';
	}"

USING INDEXES TO SPEED UP QUERIES IN POSTGRESQL:
https://niallburkley.com/blog/index-columns-for-like-in-postgres/
https://www.postgresql.org/docs/current/textsearch-indexes.html
https://www.cybertec-postgresql.com/en/postgresql-more-performance-for-like-and-ilike-statements/
https://www.yugabyte.com/blog/postgresql-like-query-performance-variations/
https://www.commandprompt.com/education/select-if-string-contains-a-substring-match-in-postgresql/
https://dev.to/____marcell/fast-fulltext-search-with-postgres-gin-index-22n5
Check trigram index:
https://stackoverflow.com/questions/54432677/why-postgresql-doesnt-use-trigram-index

*/

/* CHECKS FUNCTIONS:

// settings test:
func(s *settings.Settings) {
	fmt.Println(s.DB.Name)
},

// Database operation check function:
func(db *sqlx.DB) {
	_, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
},

// Repository operation check function.
// User Repository:
func(ctx context.Context, repo repository.Repository) {
	u := &entity.User{
		Email:    "julieta@julieta.com",
		Username: "julieta",
		Password: "password",
	}

	err := repo.SaveUser(ctx, u)
	if err != nil {
		fmt.Println(err)
	}

	user, _ := repo.GetUserByEmail(ctx, "julieta@julieta.com")

	fmt.Println(user)
},

// Link Repository:
func(ctx context.Context, repo repository.Repository) {
	l := &entity.Link{
		Url:         "another long url",
		Slug:        "1_TW3g",
		Description: "a description",
		UserID:      "a139104c-a9b4-4785-8390-840fa7b89698",
	}

	err := repo.SaveLink(ctx, l)
	if err != nil {
		fmt.Println(err)
	}

	l, err = repo.GetLink(ctx, "j6XBa9")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(l)
	ll, err := repo.GetLinks(ctx, "a139104c-a9b4-4785-8390-840fa7b89698")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ll)
},


// Service operation check function.
// Register function:
func(ctx context.Context, serv service.Service) {
	err := serv.RegisterUser(
		ctx, "enrique@enrique.com", "enrique", "123456",
	)
	if err != nil {
		panic(err)
	}
},
// Login function:
func(ctx context.Context, serv service.Service) {
	u, err := serv.LoginUser(ctx, "enrique@enrique.com", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)
},

// AddLink function:
func(ctx context.Context, serv service.Service) {
	l := &entity.Link{
		Url:    "sixth very long url",
		UserID: "a139104c-a9b4-4785-8390-840fa7b89698",
	}
	err := serv.AddLink(ctx, l)
	if err != nil {
		fmt.Println(err)
	}
},

// RecoverLinks function:
func(ctx context.Context, serv service.Service) {
	ll, err := serv.RecoverLinks(ctx, "a139104c-a9b4-4785-8390-840fa7b89698")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ll)
},

// RecoverLink function:
func(ctx context.Context, serv service.Service) {
	ml, err := serv.RecoverLink(ctx, "V_x_AH")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ml)
},

// RemoveLink function:
func(ctx context.Context, serv service.Service) {
	err := serv.RemoveLink(ctx, "7NAc20", "a139104c-a9b4-4785-8390-840fa7b89698")
	if err != nil {
		fmt.Println(err)
	}
},

*/
