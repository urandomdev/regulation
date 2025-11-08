package server

import (
	"context"
	"fmt"
	"os"

	"regulation/internal/advisor"
	"regulation/internal/config"
	"regulation/internal/ent"
	"regulation/internal/session"
	"regulation/server/services/plaid"

	"github.com/fxamacker/cbor/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"github.com/rs/zerolog"
)

var (
	encoder cbor.EncMode
	decoder cbor.DecMode
)

func init() {
	validation.ErrorTag = "cbor"

	encoder, _ = cbor.CanonicalEncOptions().EncMode()
	decoder, _ = cbor.DecOptions{}.DecMode()
}

type Server struct {
	app    *fiber.App
	config *config.Config

	db             *ent.Client
	advisorService *advisor.Service
	cache          rueidis.Client
	cacheLock      rueidislock.Locker
	sessionManager *session.Manager

	plaidClient plaid.Client
	syncService *plaid.SyncService

	logger zerolog.Logger
}

func NewServer() *Server {
	return &Server{
		app: fiber.New(fiber.Config{
			BodyLimit:          25 * 1024 * 1024, // 25 MB
			StreamRequestBody:  true,
			EnableIPValidation: true,
			ProxyHeader:        "X-Forwarded-For",
			TrustProxy:         true,
			TrustProxyConfig: fiber.TrustProxyConfig{
				Proxies: []string{
					"10.0.0.0/8",
					"172.16.0.0/12",
					"192.168.0.0/16",
					"169.254.0.0/16",
				},
			},

			RequestMethods: append(fiber.DefaultMethods, "QUERY"),

			CBOREncoder: encoder.Marshal,
			CBORDecoder: decoder.Unmarshal,
		}),
	}
}

func (s *Server) init(ctx context.Context) error {
	var err error

	s.config, err = config.LoadConfig()
	if err != nil {
		return err
	}

	level := zerolog.InfoLevel
	if s.config.Debug {
		level = zerolog.DebugLevel
	}

	s.logger = zerolog.New(os.Stderr).Level(level).With().Str("version", config.Version).Timestamp().Logger()
	s.logger = s.logger.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05 MST",
	})

	s.cache = s.config.Redis.Client()
	s.cacheLock, err = rueidislock.NewLocker(rueidislock.LockerOption{
		ClientOption:   s.config.Redis.ClientOption(),
		KeyMajority:    1,
		NoLoopTracking: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create redis lock client: %w", err)
	}

	// Initialize session manager
	s.sessionManager = session.NewManager(s.cache)

	// Initialize AI advisor using config
	apiKey := ""
	if s.config.OpenAI != nil {
		apiKey = s.config.OpenAI.APIKey
	}
	s.advisorService = advisor.NewService(apiKey)
	if apiKey == "" {
		s.logger.Warn().Msg("openai.api_key missing in config; GPT-5 budget endpoint will respond with 500")
	}

	if err = s.setupDatabase(ctx); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Initialize Plaid client
	if s.config.Plaid == nil || s.config.Plaid.UseMockClient() {
		s.logger.Warn().Str("component", "plaid").Msg("using mock plaid client; no network calls will be made")
		s.plaidClient = plaid.NewMockPlaidClient()
	} else {
		s.plaidClient, err = plaid.NewPlaidClient(s.config.Plaid)
		if err != nil {
			return fmt.Errorf("failed to create plaid client: %w", err)
		}
	}

	// Initialize sync service
	s.syncService = plaid.NewSyncService(s.plaidClient, s.db, s.config)

	s.route()

	return nil
}

func (s *Server) listen(ctx context.Context) error {
	listenAddress := ":80"
	if listenEnv, ok := os.LookupEnv("LISTEN_ADDRESS"); ok {
		listenAddress = listenEnv
	}

	return s.app.Listen(listenAddress, fiber.ListenConfig{
		GracefulContext:       ctx,
		DisableStartupMessage: !s.config.Debug,
		EnablePrintRoutes:     s.config.Debug,
	})
}

func (s *Server) shutdown() {
	if s.cacheLock != nil {
		s.logger.Info().Str("component", "redis-lock").Msg("closing redis lock client")
		s.cacheLock.Close()
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.logger.Error().Err(err).Msg("failed to close database client")
		} else {
			s.logger.Info().Str("component", "database").Msg("database client closed")
		}
	}

	if s.cache != nil {
		s.logger.Info().Str("component", "redis").Msg("closing redis client")
		s.cache.Close()
	}
}

// Run initializes and starts the server
func (s *Server) Run(ctx context.Context) error {
	if err := s.init(ctx); err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	s.logger.Info().Msg("server initialized successfully")

	defer s.shutdown()

	s.logger.Info().Msg("starting server")
	if err := s.listen(ctx); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
