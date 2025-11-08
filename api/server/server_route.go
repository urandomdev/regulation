package server

import (
	"regulation/internal/ro"
	"regulation/server/handlers/account"
	"regulation/server/middleware"
)

func (s *Server) route() {
	s.app.Use(middleware.NewRequestInfoMiddleware())
	s.app.Use(middleware.NewLoggerMiddleware(s.logger))

	accountGroup := s.app.Group("/account")
	{
		handler := account.New(s.db, s.sessionManager)
		auth := middleware.NewAuth(s.db, s.sessionManager)

		accountGroup.Post("/signup", ro.WrapHandler2(handler.Signup))
		accountGroup.Post("/login", ro.WrapHandler2(handler.Login))

		// Protected routes (authentication required)
		accountGroup.Post("/logout", auth.Handle, ro.WrapHandler4(handler.Logout))
		accountGroup.Get("/me", auth.Handle, ro.WrapHandler3(handler.Me))
	}

}
