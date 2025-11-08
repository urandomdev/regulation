package server

import (
	"regulation/internal/ro"
	"regulation/server/handlers/account"
)

func (s *Server) route() {
	handler := account.New(s.db)

	s.app.Post("/account/signup", ro.WrapHandler2(handler.Signup))
}
