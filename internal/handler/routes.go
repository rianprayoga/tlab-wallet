package handler

import (
	"net/http"
	"tlab-wallet/internal/config"
	"tlab-wallet/internal/repository"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
)

type HttpHandler struct {
	Db  repository.Repo
	Log *logrus.Logger
	Jwt *jwtauth.JWTAuth
}

func (h *HttpHandler) Routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	mux.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)
	})

	mux.Group(func(r chi.Router) {
		r.Use(config.Verify(h.Jwt))
		r.Use(config.Authenticator(h.Jwt))

		r.Get("/api/users/profile", h.GetProfile)

		r.Route("/api/wallets", func(cr chi.Router) {
			cr.Get("/balance", h.GetWallet)
			cr.Post("/topup", h.TopUpWallet)
		})

		r.Route("/api/transactions", func(cr chi.Router) {
			cr.Get("/history", h.GetTransactions)
			cr.Post("/transfer", h.Transfer)
		})

	})

	return mux
}
