// @title           Gachita API
// @version         1.0
// @description     gachita API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"gachita-api/internal/auth"
	"gachita-api/internal/config"
	"gachita-api/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"

	httpadapter "gachita-api/internal/adapter/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(1)
	}

	defer pool.Close()
	queries := db.New(pool)
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTExpiresIn)
	handler := httpadapter.NewRouter(pool, queries, tokens)
	fmt.Println("Server is running on", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, handler); err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
