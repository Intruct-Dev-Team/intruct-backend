package application

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/config"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/repository"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/services/course"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/services/language"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/services/lesson"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/services/n8n"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/services/user"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/jwt"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/migrator"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/db/postgres"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/object/supabase"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Application struct {
	db     *sqlx.DB
	server *rest.Server
	log    *zap.Logger
}

type Services struct {
	jwt.JWTService
	n8n.N8NService
	user.UserService
	language.LanguageService
	course.CourseService
	lesson.LessonService
}

func NewServices(
	jwtService *jwt.JWTService,
	n8nService *n8n.N8NService,
	userService *user.UserService,
	languageService *language.LanguageService,
	courseService *course.CourseService,
	lessonService *lesson.LessonService,
) *Services {
	return &Services{
		JWTService:      *jwtService,
		N8NService:      *n8nService,
		UserService:     *userService,
		LanguageService: *languageService,
		CourseService:   *courseService,
		LessonService:   *lessonService,
	}
}

func New(ctx context.Context, config *config.Config, log *zap.Logger) (*Application, error) {
	var err error

	// database connection
	database, err := postgres.New(ctx, &config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	log.Info("connected to database successfully")

	repo := repository.New(database)

	if config.IsInitDb {
		err = migrator.RunMigrations(&config.Migrator)
		if err != nil {
			return nil, err
		}

		log.Info("up migrations successfully")
	}

	/*----------------------------------------------------------*/

	// services
	jwtService := jwt.NewService(config.JwtSecretKey, jwt.WithIssuer("learn-share-backend"))
	n8nService := n8n.NewService(config.N8NApiRoute, log)
	fileService := supabase.NewService(&config.ObjectStorage)

	userService := user.NewService(repo, fileService)
	languageService := language.NewService(repo)
	courseService := course.NewService(repo)
	lessonService := lesson.NewService(repo)

	services := NewServices(
		jwtService,
		n8nService,
		userService,
		languageService,
		courseService,
		lessonService,
	)

	restServer := rest.NewServer(services, config.Server, log)

	return &Application{
		db:     database,
		server: restServer,
		log:    log,
	}, nil
}

// Run start application.
func (app *Application) Run() error {
	err := app.server.Start()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully stop application.
func (app *Application) Shutdown(ctx context.Context) error {
	app.log.Info("shutting down application...")

	go func() {
		<-ctx.Done()

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			app.log.Error("graceful shutdown timed out... forcing exit")
			os.Exit(1)
		}
	}()

	if err := app.server.GracefulStop(ctx); err != nil {
		return err
	}

	if err := app.db.Close(); err != nil {
		app.log.Error("failed to close database", zap.Error(err))
	}

	return nil
}
