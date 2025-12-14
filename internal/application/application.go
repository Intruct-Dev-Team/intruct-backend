package application

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/config"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/repository"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/services/user"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/jwt"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/migrator"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/db/postgres"
	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/object/s3"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Application struct {
	db     *sqlx.DB
	server *rest.Server
	log    *zap.Logger
}

type Services struct {
	user.UserService
	jwt.JWTService
	// 	kratos.KratosService
	// 	teacher.TeacherService
	// 	schedule.ScheduleService
	// 	review.ReviewService
	// 	lesson.LessonService
	// 	image.ImageService
	// 	category.CategoryService
	// 	skill.SkillService
	// 	complaint.ComplaintService
	// 	common.CommonService
}

func NewServices(
	userService *user.UserService,
	jwtService *jwt.JWTService,
	// kratosService *kratos.KratosService,
	// teacherService *teacher.TeacherService,
	// scheduleService *schedule.ScheduleService,
	// reviewService *review.ReviewService,
	// lessonService *lesson.LessonService,
	// imageService *image.ImageService,
	// categoryService *category.CategoryService,
	// skillService *skill.SkillService,
	// complaintService *complaint.ComplaintService,
	// commonService *common.CommonService,
) *Services {
	return &Services{
		UserService: *userService,
		JWTService:  *jwtService,
		// 		KratosService:    *kratosService,
		// 		TeacherService:   *teacherService,
		// 		ScheduleService:  *scheduleService,
		// 		ReviewService:    *reviewService,
		// 		LessonService:    *lessonService,
		// 		ImageService:     *imageService,
		// 		CategoryService:  *categoryService,
		// 		SkillService:     *skillService,
		// 		ComplaintService: *complaintService,
		// 		CommonService:    *commonService,
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

	s3Client, err := s3.NewClient(&config.S3)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to S3 object storage (supabase): %w", err)
	}

	log.Info("connected to S3 object storage successfully")

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
	// kratosService := kratos.New(config.Kratos)
	jwtService := jwt.NewService(config.JwtSecretKey, jwt.WithIssuer("learn-share-backend"))
	s3Service := s3.NewService(s3Client, config.S3.Bucket, config.S3.Host)
	// liveKitService := livekit.NewService(config.LiveKit)
	// commonService := common.NewService(repo)

	userService := user.NewService(repo, s3Service)
	// teacherService := teacher.NewService(repo)
	// scheduleService := schedule.NewService(repo)
	// reviewService := review.NewService(repo)
	// lessonService := lesson.NewService(repo, liveKitService)
	// imageService := image.NewService(minioService)
	// categoryService := category.NewService(repo)
	// skillService := skill.NewService(repo)
	// complaintService := complaint.NewService(repo)

	services := NewServices(
		userService,
		jwtService,
	// 	kratosService,
	// 	jwtService,
	// 	teacherService,
	// 	scheduleService,
	// 	reviewService,
	// 	lessonService,
	// 	imageService,
	// 	categoryService,
	// 	skillService,
	// 	complaintService,
	// 	commonService,
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

	// if err := app.db.Close(); err != nil {
	// 	app.log.Error("failed to close database", zap.Error(err))
	// }

	return nil
}
