package httpadapter

import (
	"fmt"
	"net/http"

	"gachita-api/internal/auth"
	"gachita-api/internal/db"

	_ "gachita-api/docs"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	tokens  *auth.TokenService
}

func NewRouter(pool *pgxpool.Pool, queries *db.Queries, tokens *auth.TokenService) http.Handler {
	r := &Router{pool: pool, queries: queries, tokens: tokens}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", r.health)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// 인증 불필요
	mux.HandleFunc("POST /api/auth/signup", r.signUp)
	mux.HandleFunc("POST /api/auth/login", r.login)

	// 인증 필요
	mux.Handle("GET /api/me", r.authMiddleware(http.HandlerFunc(r.getMe)))
	mux.Handle("GET /api/me/queue", r.authMiddleware(http.HandlerFunc(r.listMyQueueEntries)))
	mux.Handle("GET /api/rooms", r.authMiddleware(http.HandlerFunc(r.listRooms)))
	mux.Handle("POST /api/rooms", r.authMiddleware(http.HandlerFunc(r.createRoom)))
	mux.Handle("POST /api/rooms/join", r.authMiddleware(http.HandlerFunc(r.joinRoom)))
	mux.Handle("GET /api/rooms/{id}", r.authMiddleware(http.HandlerFunc(r.getRoom)))
	mux.Handle("POST /api/rooms/{id}/stops", r.authMiddleware(http.HandlerFunc(r.createHubStop)))
	mux.Handle("GET /api/rooms/{id}/stops", r.authMiddleware(http.HandlerFunc(r.listHubStops)))
	mux.Handle("PUT /api/rooms/{id}/stops/{stopId}", r.authMiddleware(http.HandlerFunc(r.updateHubStop)))
	mux.Handle("DELETE /api/rooms/{id}/stops/{stopId}", r.authMiddleware(http.HandlerFunc(r.deleteHubStop)))
	mux.Handle("POST /api/rooms/{id}/queue", r.authMiddleware(http.HandlerFunc(r.createQueueEntry)))
	mux.Handle("GET /api/rooms/{id}/queue", r.authMiddleware(http.HandlerFunc(r.listQueueEntries)))
	mux.Handle("GET /api/rooms/{id}/queue/{entryId}", r.authMiddleware(http.HandlerFunc(r.getQueueEntry)))
	mux.Handle("DELETE /api/rooms/{id}/queue/{entryId}", r.authMiddleware(http.HandlerFunc(r.cancelQueueEntry)))
	mux.Handle("GET /api/rooms/{id}/matches", r.authMiddleware(http.HandlerFunc(r.listMatches)))
	mux.Handle("GET /api/rooms/{id}/matches/{matchId}", r.authMiddleware(http.HandlerFunc(r.getMatch)))
	mux.Handle("GET /api/notifications", r.authMiddleware(http.HandlerFunc(r.listNotifications)))
	mux.Handle("PUT /api/notifications/read-all", r.authMiddleware(http.HandlerFunc(r.readAllNotifications)))
	mux.Handle("PUT /api/notifications/{notificationId}/read", r.authMiddleware(http.HandlerFunc(r.readNotification)))
	return corsMiddleware(mux)
}

func (r *Router) health(w http.ResponseWriter, req *http.Request) {
	ok, err := r.queries.PingTest(req.Context())
	if err != nil {
		http.Error(w, "db ping failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "db ping ok: %d", ok)
}
