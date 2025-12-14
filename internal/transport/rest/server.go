package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/handlers"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/httputils"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest/middlewares"
	"github.com/go-chi/chi/v5"

	// httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

const (
	defaultHTTPServerWriteTimeout = time.Second * 15
	defaultHTTPServerReadTimeout  = time.Second * 15

	apiRoute = "/api"
)

type Services interface {
	handlers.Services

	// middlewares.SessionValidator
	middlewares.UserService    // for auth
	middlewares.TokenValidator // for auth
}

type Config struct {
	Port int `env:"SERVER_PORT" env-required:"true"`
}

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(services Services, config Config, log *zap.Logger) *Server {
	router := chi.NewRouter()

	router.Use(middlewares.LoggerMiddleware(log.Named("log_middleware")))
	router.Use(middlewares.CorsMiddleware)

	var tokenValidator middlewares.TokenValidator = services
	var userService middlewares.UserService = services
	jwtAuthMiddleware := middlewares.JWTAuthMiddleware(tokenValidator, log.Named("jwt_middleware"))
	systemAuthMiddleware := middlewares.SystemAuthMiddleware(userService, log.Named("sys_auth_middlware"))

	handler := handlers.New(services, log)

	// root router
	apiRouter := chi.NewRouter()

	apiRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httputils.RespondWith200(w, struct{ status string }{status: "ok"}, log)
	})

	// all routes
	handler.SetupRoutes(apiRouter, jwtAuthMiddleware, systemAuthMiddleware)

	router.Mount(apiRoute, apiRouter)

	// add swagger endpoint
	// router.Get("/swagger/*", httpSwagger.Handler(
	// 	httpSwagger.URL("./swagger/doc.json"), // URL to JSON docs file
	// ))

	return &Server{
		server: &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf(":%d", config.Port),
			WriteTimeout: defaultHTTPServerWriteTimeout,
			ReadTimeout:  defaultHTTPServerReadTimeout,
		},
		logger: log,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting Rest server", zap.String("address", s.server.Addr))

	return s.server.ListenAndServe()
}

// GracefulStop right server stopping
func (s *Server) GracefulStop(ctx context.Context) error {
	// create context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s.logger.Info("shutting down Rest server", zap.String("address", s.server.Addr))

	// stop server
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		err = fmt.Errorf("failed to shutdown rest server: %w", err)

		return err
	}

	s.logger.Info("rest server stopped")

	return nil
}
