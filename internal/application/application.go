package application

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Intruct-Dev-Team/intruct-backend/internal/config"
	"github.com/Intruct-Dev-Team/intruct-backend/internal/transport/rest"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"os"

// 	"github.com/Intruct-Dev-Team/intruct-backend/pkg/jwt"
// 	"github.com/Intruct-Dev-Team/intruct-backend/pkg/migrator"
// 	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/db/postgres"
// 	"github.com/Intruct-Dev-Team/intruct-backend/pkg/storage/object/minio"
// 	"github.com/LearnShareApp/learn-share-backend/pkg/livekit"
// 	"github.com/jmoiron/sqlx"
// 	"go.uber.org/zap"
// 	"honnef.co/go/tools/config"
// )

type Application struct {
	db     *sqlx.DB
	server *rest.Server
	log    *zap.Logger
}

// type Services struct {
// 	kratos.KratosService
// 	jwt.JWTService
// 	user.UserService
// 	teacher.TeacherService
// 	schedule.ScheduleService
// 	review.ReviewService
// 	lesson.LessonService
// 	image.ImageService
// 	category.CategoryService
// 	skill.SkillService
// 	complaint.ComplaintService
// 	common.CommonService
// }

// func NewServices(
// 	kratosService *kratos.KratosService,
// 	jwtService *jwt.JWTService,
// 	userService *user.UserService,
// 	teacherService *teacher.TeacherService,
// 	scheduleService *schedule.ScheduleService,
// 	reviewService *review.ReviewService,
// 	lessonService *lesson.LessonService,
// 	imageService *image.ImageService,
// 	categoryService *category.CategoryService,
// 	skillService *skill.SkillService,
// 	complaintService *complaint.ComplaintService,
// 	commonService *common.CommonService,
// ) *Services {
// 	return &Services{
// 		KratosService:    *kratosService,
// 		JWTService:       *jwtService,
// 		UserService:      *userService,
// 		TeacherService:   *teacherService,
// 		ScheduleService:  *scheduleService,
// 		ReviewService:    *reviewService,
// 		LessonService:    *lessonService,
// 		ImageService:     *imageService,
// 		CategoryService:  *categoryService,
// 		SkillService:     *skillService,
// 		ComplaintService: *complaintService,
// 		CommonService:    *commonService,
// 	}
// }

func New(ctx context.Context, config *config.Config, log *zap.Logger) (*Application, error) {
	// database connection
	// database, err := postgres.New(ctx, &config.DB)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	// }

	// log.Info("connected to database successfully")

	// minioClient, err := minio.NewClient(&config.Minio)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to minio: %w", err)
	// }

	// if err = minio.CreateBucket(ctx, minioClient, config.Minio.Bucket); err != nil {
	// 	return nil, fmt.Errorf("failed to create minio bucket: %w", err)
	// }

	// log.Info("connected to minio successfully")

	// repo := repository.New(database)

	// if config.IsInitDb {
	// 	err = migrator.RunMigrations(&config.Migrator)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	log.Info("up migrations successfully")
	// }

	/*----------------------------------------------------------*/

	// services
	// kratosService := kratos.New(config.Kratos)
	// jwtService := jwt.NewService(config.JwtSecretKey, jwt.WithIssuer("learn-share-backend"))
	// liveKitService := livekit.NewService(config.LiveKit)
	// minioService := minio.NewService(minioClient, config.Minio.Bucket)
	// commonService := common.NewService(repo)

	// userService := user.NewService(repo, minioService)
	// teacherService := teacher.NewService(repo)
	// scheduleService := schedule.NewService(repo)
	// reviewService := review.NewService(repo)
	// lessonService := lesson.NewService(repo, liveKitService)
	// imageService := image.NewService(minioService)
	// categoryService := category.NewService(repo)
	// skillService := skill.NewService(repo)
	// complaintService := complaint.NewService(repo)

	// services := NewServices(
	// 	kratosService,
	// 	jwtService,
	// 	userService,
	// 	teacherService,
	// 	scheduleService,
	// 	reviewService,
	// 	lessonService,
	// 	imageService,
	// 	categoryService,
	// 	skillService,
	// 	complaintService,
	// 	commonService,
	// )

	restServer := rest.NewServer(struct{}{}, config.Server, log)

	return &Application{
		// db:     database,
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
