package httpadapter

import (
	"fmt"
	"net/http"

	"gachita-api/internal/auth"
	"gachita-api/internal/db"

	_ "gachita-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	queries *db.Queries
	tokens  *auth.TokenService
}

func NewRouter(queries *db.Queries, tokens *auth.TokenService) http.Handler {
	r := &Router{queries: queries, tokens: tokens}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", r.health)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /api/auth/signup", r.signUp)
	mux.HandleFunc("POST /api/auth/login", r.login)
	return mux
}

func (r *Router) health(w http.ResponseWriter, req *http.Request) {
	ok, err := r.queries.PingTest(req.Context())
	if err != nil {
		http.Error(w, "db ping failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "db ping ok: %d", ok)
}
