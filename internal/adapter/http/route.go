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

	// 인증 불필요
	mux.HandleFunc("POST /api/auth/signup", r.signUp)
	mux.HandleFunc("POST /api/auth/login", r.login)

	// 인증 필요
	mux.Handle("GET /api/me", r.authMiddleware(http.HandlerFunc(r.getMe)))
	mux.Handle("POST /api/rooms", r.authMiddleware(http.HandlerFunc(r.createRoom)))
	mux.Handle("POST /api/rooms/join", r.authMiddleware(http.HandlerFunc(r.joinRoom)))
	mux.Handle("GET /api/rooms/{id}", r.authMiddleware(http.HandlerFunc(r.getRoom)))
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
