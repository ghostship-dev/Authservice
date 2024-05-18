package core

import (
	"fmt"

	"github.com/ghostship-dev/authservice/core/handlers"
	"github.com/ghostship-dev/authservice/core/router"
	_ "github.com/joho/godotenv/autoload"
)

func RunService() {
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

	fmt.Println("Starting service...")

	err := apiV1Router.ListenAndServe(":8080")
	if err != nil {
		panic(err)
	}
}
