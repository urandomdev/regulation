package server

import (
	"regulation/internal/ro"
	"regulation/server/handlers/account"
	advisorhandler "regulation/server/handlers/advisor"
	"regulation/server/handlers/financial"
	"regulation/server/handlers/notification"
	"regulation/server/handlers/plaid"
	"regulation/server/middleware"
)

func (s *Server) route() {
	s.app.Use(middleware.NewRequestInfoMiddleware())
	s.app.Use(middleware.NewLoggerMiddleware(s.logger))

	auth := middleware.NewAuth(s.db, s.sessionManager)

	accountGroup := s.app.Group("/account")
	{
		handler := account.New(s.db, s.sessionManager)

		accountGroup.Post("/signup", ro.WrapHandler2(handler.Signup))
		accountGroup.Post("/login", ro.WrapHandler2(handler.Login))

		// Protected routes (authentication required)
		accountGroup.Post("/logout", auth.Handle, ro.WrapHandler4(handler.Logout))
		accountGroup.Get("/me", auth.Handle, ro.WrapHandler3(handler.Me))
	}

	// Plaid onboarding routes - for linking bank accounts
	plaidGroup := s.app.Group("/plaid")
	{
		handler := plaid.New(s.db, s.sessionManager, s.plaidClient, s.syncService)

		// All Plaid routes require authentication
		plaidGroup.Post("/create-link-token", auth.Handle, ro.WrapHandler3(handler.CreateLinkToken))
		plaidGroup.Post("/exchange-token", auth.Handle, ro.WrapHandler2(handler.ExchangeToken))
		plaidGroup.Post("/sync-transactions", auth.Handle, ro.WrapHandler2(handler.SyncTransactions))
		plaidGroup.Delete("/accounts/:id", auth.Handle, ro.WrapHandler3(handler.DisconnectAccount))
	}

	// Financial dashboard routes - for viewing account data
	financialGroup := s.app.Group("/financial")
	{
		handler := financial.New(s.db, s.sessionManager)

		// All financial routes require authentication
		financialGroup.Get("/accounts", auth.Handle, ro.WrapHandler3(handler.GetAccounts))
		financialGroup.Post("/transactions", auth.Handle, ro.WrapHandler(handler.GetTransactions))
		financialGroup.Post("/accounts/:id/transactions", auth.Handle, ro.WrapHandler(handler.GetAccountTransactions))
		financialGroup.Post("/cashflow", auth.Handle, ro.WrapHandler(handler.GetCashflow))
	}

	advisorGroup := s.app.Group("/advisor")
	{
		handler := advisorhandler.New(s.db, s.advisorService)
		advisorGroup.Post("/budget-plan/test", ro.WrapHandler(handler.BudgetPlan))
		advisorGroup.Post("/budget-plan/history", ro.WrapHandler(handler.BudgetPlanFromHistory))
	}

	notificationGroup := s.app.Group("/notification")
	{
		handler := notification.New(s.config, s.db)

		// Public endpoint to get VAPID public key
		notificationGroup.Get("/vapid", ro.WrapHandler3(handler.GetVAPIDPublicKey))

		// Protected endpoints (authentication required)
		notificationGroup.Post("/subscribe", auth.Handle, ro.WrapHandler2(handler.Subscribe))
		notificationGroup.Delete("/subscribe", auth.Handle, ro.WrapHandler2(handler.Unsubscribe))
	}
}
