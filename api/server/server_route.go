package server

import (
	"regulation/internal/ro"
	"regulation/server/handlers/account"
	"regulation/server/middleware"
)

func (s *Server) route() {
	handler := account.New(s.db, s.sessionManager)
	auth := middleware.NewAuth(s.db, s.sessionManager)

	// Public routes (no authentication required)
	s.app.Post("/account/signup", ro.WrapHandler2(handler.Signup))
	s.app.Post("/account/login", ro.WrapHandler2(handler.Login))

	// Protected routes (authentication required)
	s.app.Post("/account/logout", auth.Handle, handler.Logout)
	s.app.Get("/account/me", auth.Handle, handler.Me)
}
