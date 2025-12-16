package main

import (
	"log"
	"net/http"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/handler"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/config"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/database"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database connected")

	gormDB := db.GetDB()

	log.Println("🔄 Running migrations...")
	if err := database.RunMigrations(gormDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("✅ Migrations completed")

	// Repositories
	userRepo := repository.NewPostgresUserRepository(gormDB)
	courseRepo := repository.NewPostgresCourseRepository(gormDB)
	enrollmentRepo := repository.NewPostgresEnrollmentRepository(gormDB)
	moduleRepo := repository.NewPostgresModuleRepository(gormDB)
	lessonRepo := repository.NewPostgresLessonRepository(gormDB)

	tokenProvider := auth.NewJWTProvider(&cfg.JWT)

	// Use Cases
	userUseCase := usecase.NewUserUseCase(userRepo, tokenProvider)
	courseUseCase := usecase.NewCourseUseCase(courseRepo, userRepo)
	enrollmentUseCase := usecase.NewEnrollmentUseCase(enrollmentRepo, courseRepo, userRepo)
	moduleUseCase := usecase.NewModuleUseCase(moduleRepo, courseRepo)
	lessonUseCase := usecase.NewLessonUseCase(lessonRepo, moduleRepo, courseRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userUseCase)
	courseHandler := handler.NewCourseHandler(courseUseCase)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentUseCase)
	moduleHandler := handler.NewModuleHandler(moduleUseCase)
	lessonHandler := handler.NewLessonHandler(lessonUseCase)

	router := SetupRoutes(userHandler, courseHandler, enrollmentHandler, moduleHandler, lessonHandler, tokenProvider)

	router = middleware.LoggerMiddleware(router)
	router = middleware.CORSMiddleware(router)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
