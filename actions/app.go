package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"github.com/devrelcollective/xela/models"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr"

	"github.com/markbates/goth/gothic"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_xela_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		app.GET("/", HomeHandler)

		auth := app.Group("/auth")
		app.Use(SetCurrentUser)
		app.Use(Authorize)
		app.Middleware.Skip(Authorize, HomeHandler)
		bah := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("/{provider}", bah)
		auth.DELETE("", AuthDestroy)
		auth.Middleware.Skip(Authorize, bah, AuthCallback)
		auth.GET("/{provider}/callback", AuthCallback)
		// TODO Put all these resources into a group like auth is
		dr := DutoniansResource{&buffalo.BaseResource{}}
		dg := app.Resource("/dutonians", dr)
		dg.Middleware.Skip(Authorize, dr.List, dr.Show)
		er := EventsResource{&buffalo.BaseResource{}}
		eg := app.Resource("/events", er)
		eg.Middleware.Skip(Authorize, er.List, er.Show)
		ar := AbstractsResource{&buffalo.BaseResource{}}
		ag := app.Resource("/abstracts", ar)
		ag.Middleware.Skip(Authorize, ar.List, ar.Show)
		pr := ProposalsResource{&buffalo.BaseResource{}}
		pg := app.Resource("/proposals", pr)
		pg.Middleware.Skip(Authorize, pr.List, pr.Show)
		sr := SponsorshipsResource{&buffalo.BaseResource{}}
		sg := app.Resource("/sponsorships", sr)
		sg.Middleware.Skip(Authorize, sr.List, sr.Show)
		lr := LinksResource{&buffalo.BaseResource{}}
		lg := app.Resource("/links", lr)
		lg.Middleware.Skip(Authorize, lr.List, lr.Show)
		// app.Resource("/dutonians", DutoniansResource{})
		// app.Resource("/events", EventsResource{})
		// app.Resource("/abstracts", AbstractsResource{})
		// app.Resource("/proposals", ProposalsResource{})
		// app.Resource("/sponsorships", SponsorshipsResource{})
		// app.Resource("/links", LinksResource{})
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
