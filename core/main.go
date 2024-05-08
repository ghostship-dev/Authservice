package core

import (
	"fmt"
	"github.com/ghostship-dev/authservice/core/handlers"
	"github.com/ghostship-dev/authservice/core/router"
)

func RunService() {
	apiV1Router := router.New().Group("/api/v1")

	apiV1Router.Post("/login", handlers.LoginHandler)
	apiV1Router.Post("/register", handlers.RegisterHandler)

	fmt.Println("Starting service...")

	err := apiV1Router.ListenAndServe(":8080")
	if err != nil {
		panic(err)
	}
}
