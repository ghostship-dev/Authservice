package handlers

import (
	"github.com/ghostship-dev/authservice/core/responses"
	"net/http"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	responses.NewTextResponse(w, 200, "OK")
}
