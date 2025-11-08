package server

import (
	"regulation/internal/ro"
	"regulation/server/handlers/account"
	"regulation/server/handlers/financial"
	"regulation/server/handlers/notification"
	"regulation/server/handlers/plaid"
	"regulation/server/handlers/recommendation"
	"regulation/server/handlers/rule"
	"regulation/server/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	recoverHandler "github.com/gofiber/fiber/v3/middleware/recover"
)

func (s *Server) route() {
	s.app.Use(recoverHandler.New(recoverHandler.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx fiber.Ctx, e any) {
			s.logger.Error().Any("panic", e).Msg("Panic recovered in HTTP handler")
		},
	}))
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     s.config.CORS.AllowedOrigins,
		AllowCredentials: true,
	}))

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
		generatorHandler := plaid.NewTransactionGeneratorHandler(handler)

		// All Plaid routes require authentication
		plaidGroup.Post("/create-link-token", auth.Handle, ro.WrapHandler3(handler.CreateLinkToken))
		plaidGroup.Post("/exchange-token", auth.Handle, ro.WrapHandler2(handler.ExchangeToken))
		plaidGroup.Post("/sync-transactions", auth.Handle, ro.WrapHandler2(handler.SyncTransactions))
		plaidGroup.Delete("/accounts/:id", auth.Handle, ro.WrapHandler3(handler.DisconnectAccount))

		// Transaction generator endpoints (sandbox only)
		plaidGroup.Post("/generator/start", auth.Handle, ro.WrapHandler2(generatorHandler.StartGenerator))
		plaidGroup.Post("/generator/stop", auth.Handle, ro.WrapHandler2(generatorHandler.StopGenerator))
		plaidGroup.Get("/generator/status/:item_id", auth.Handle, ro.WrapHandler4(generatorHandler.GetGeneratorStatus))
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

	notificationGroup := s.app.Group("/notification")
	{
		handler := notification.New(s.config, s.db)

		// Public endpoint to get VAPID public key
		notificationGroup.Get("/vapid", ro.WrapHandler3(handler.GetVAPIDPublicKey))

		// Protected endpoints (authentication required)
		notificationGroup.Post("/subscribe", auth.Handle, ro.WrapHandler2(handler.Subscribe))
		notificationGroup.Delete("/subscribe", auth.Handle, ro.WrapHandler2(handler.Unsubscribe))
	}

	// Rule management routes - for creating and managing savings rules
	ruleGroup := s.app.Group("/rules")
	{
		ruleHandler := rule.New(s.db, s.sessionManager)

		// All rule routes require authentication
		ruleGroup.Post("/", auth.Handle, ro.WrapHandler(ruleHandler.CreateRule))
		ruleGroup.Get("/", auth.Handle, ro.WrapHandler3(ruleHandler.ListRules))
		ruleGroup.Get("/:id", auth.Handle, ro.WrapHandler3(ruleHandler.GetRule))
		ruleGroup.Patch("/:id", auth.Handle, ro.WrapHandler(ruleHandler.UpdateRule))
		ruleGroup.Delete("/:id", auth.Handle, ro.WrapHandler4(ruleHandler.DeleteRule))
		ruleGroup.Patch("/:id/toggle", auth.Handle, ro.WrapHandler3(ruleHandler.ToggleRule))
		ruleGroup.Get("/:id/executions", auth.Handle, ro.WrapHandler3(ruleHandler.GetRuleExecutions))
	}

	// Rule recommendation routes - AI-powered rule suggestions
	// Note: To create a rule from a recommendation, clients should use POST /rules
	recommendationGroup := s.app.Group("/recommendations")
	{
		recommendationHandler := recommendation.New(s.db, s.sessionManager, s.suggestionService)

		// All recommendation routes require authentication
		recommendationGroup.Get("/", auth.Handle, ro.WrapHandler3(recommendationHandler.GetRecommendations))
	}
}
