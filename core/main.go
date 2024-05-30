package core

import (
	"fmt"

	"github.com/ghostship-dev/authservice/core/config"
	"github.com/ghostship-dev/authservice/core/database"
	"github.com/ghostship-dev/authservice/core/handlers"
	"github.com/ghostship-dev/authservice/core/router"
	_ "github.com/joho/godotenv/autoload"
)

func RunService(c *config.Config) {
	database.Connection = database.ConnectToSelectedDBDriver(c)

	apiV1Router := router.New().Group("/api/v1")

	// Account management
	apiV1Router.Post("/login", handlers.LoginHandler)
	apiV1Router.Post("/register", handlers.RegisterHandler)

	// Time-Based One-Time Password management
	apiV1Router.Post("/otp", handlers.AccountOTP)

	// OAuth2 Client-Application management
	apiV1Router.Post("/oauth/application", handlers.NewOAuthApplication)
	apiV1Router.Patch("/oauth/application", handlers.UpdateOAuthApplicationKeyValue)
	apiV1Router.Delete("/oauth/application", handlers.DeleteOAuthClientApplication)

	// OAuth2 Implementation
	apiV1Router.Post("/oauth/token/introspect", handlers.IntrospectOAuthToken)
	apiV1Router.Get("/oauth/authorize", handlers.AuthorizeOAuthApplication)
	apiV1Router.Post("/oauth/token", handlers.OAuthTokenEndpoint)
	apiV1Router.Post("/oauth/token/revoke", handlers.RevokeOAuthToken)

	fmt.Println(fmt.Sprintf("Running Service on: %s:%d", c.Hostname, c.Port))

	err := apiV1Router.ListenAndServe(fmt.Sprintf("%s:%d", c.Hostname, c.Port))
	if err != nil {
		panic(err)
	}
}
