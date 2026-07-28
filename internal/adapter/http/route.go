package httpadapter

import (
	"fmt"
	"net/http"

	"gachita-api/internal/db"
)

type Router struct {
	queries *db.Queries
}

func NewRouter(queries *db.Queries) http.Handler {
	r := &Router{queries: queries}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", r.health)
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
