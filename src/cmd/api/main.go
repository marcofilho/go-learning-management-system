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

// @title Learning Management System API
// @version 1.0
// @description REST API for a Learning Management System with courses, modules, lessons, and enrollments
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@lms.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

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

	userRepo := repository.NewPostgresUserRepository(gormDB)
	courseRepo := repository.NewPostgresCourseRepository(gormDB)
	enrollmentRepo := repository.NewPostgresEnrollmentRepository(gormDB)
	moduleRepo := repository.NewPostgresModuleRepository(gormDB)
	lessonRepo := repository.NewPostgresLessonRepository(gormDB)

	tokenProvider := auth.NewJWTProvider(&cfg.JWT)

	userUseCase := usecase.NewUserUseCase(userRepo, tokenProvider)
	courseUseCase := usecase.NewCourseUseCase(courseRepo, userRepo)
	enrollmentUseCase := usecase.NewEnrollmentUseCase(enrollmentRepo, courseRepo, userRepo)
	moduleUseCase := usecase.NewModuleUseCase(moduleRepo, courseRepo)
	lessonUseCase := usecase.NewLessonUseCase(lessonRepo, moduleRepo, courseRepo)

	userHandler := handler.NewUserHandler(userUseCase)
	courseHandler := handler.NewCourseHandler(courseUseCase)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentUseCase)
	moduleHandler := handler.NewModuleHandler(moduleUseCase)
	lessonHandler := handler.NewLessonHandler(lessonUseCase)

	router := SetupRoutes(userHandler, courseHandler, enrollmentHandler, moduleHandler, lessonHandler, tokenProvider, courseRepo, enrollmentRepo, moduleRepo)

	router = middleware.LoggerMiddleware(router)
	router = middleware.CORSMiddleware(router)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
