package main

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/docs"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/handler"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
)

func SetupRoutes(
	userHandler *handler.UserHandler,
	courseHandler *handler.CourseHandler,
	enrollmentHandler *handler.EnrollmentHandler,
	moduleHandler *handler.ModuleHandler,
	lessonHandler *handler.LessonHandler,
	tokenProvider auth.TokenProvider,
) http.Handler {
	r := mux.NewRouter()

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := r.PathPrefix("/api").Subrouter()

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// Authentication routes (public)
	api.HandleFunc("/auth/register", userHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", userHandler.Login).Methods(http.MethodPost)

	// User routes (protected)
	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware(tokenProvider))
	users.HandleFunc("", userHandler.ListUsers).Methods(http.MethodGet)
	users.HandleFunc("/{id}", userHandler.GetUserByID).Methods(http.MethodGet)
	users.HandleFunc("/{id}", userHandler.UpdateUser).Methods(http.MethodPut)
	users.HandleFunc("/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)

	// Course routes
	courses := api.PathPrefix("/courses").Subrouter()
	courses.Use(middleware.AuthMiddleware(tokenProvider))
	courses.HandleFunc("", courseHandler.ListCourses).Methods(http.MethodGet)
	courses.HandleFunc("", courseHandler.CreateCourse).Methods(http.MethodPost)
	courses.HandleFunc("/{id}", courseHandler.GetCourse).Methods(http.MethodGet)
	courses.HandleFunc("/{id}", courseHandler.UpdateCourse).Methods(http.MethodPut)
	courses.HandleFunc("/{id}", courseHandler.DeleteCourse).Methods(http.MethodDelete)

	// Course enrollment routes
	courses.HandleFunc("/{id}/enroll", enrollmentHandler.EnrollInCourse).Methods(http.MethodPost)
	courses.HandleFunc("/{id}/students", enrollmentHandler.GetCourseStudents).Methods(http.MethodGet)

	// Course modules routes
	courses.HandleFunc("/{courseId}/modules", moduleHandler.CreateModule).Methods(http.MethodPost)
	courses.HandleFunc("/{courseId}/modules", moduleHandler.GetCourseModules).Methods(http.MethodGet)

	// Module routes
	modules := api.PathPrefix("/modules").Subrouter()
	modules.Use(middleware.AuthMiddleware(tokenProvider))
	modules.HandleFunc("/{id}", moduleHandler.UpdateModule).Methods(http.MethodPut)
	modules.HandleFunc("/{id}", moduleHandler.DeleteModule).Methods(http.MethodDelete)

	// Module lesson routes
	modules.HandleFunc("/{moduleId}/lessons", lessonHandler.CreateLesson).Methods(http.MethodPost)
	modules.HandleFunc("/{moduleId}/lessons", lessonHandler.GetModuleLessons).Methods(http.MethodGet)

	// Lesson routes
	lessons := api.PathPrefix("/lessons").Subrouter()
	lessons.Use(middleware.AuthMiddleware(tokenProvider))
	lessons.HandleFunc("/{lessonId}/version", lessonHandler.CreateLessonVersion).Methods(http.MethodPost)
	lessons.HandleFunc("/{lessonId}/all-versions", lessonHandler.GetAllLessonVersions).Methods(http.MethodGet)
	lessons.HandleFunc("/{lessonId}", lessonHandler.DeleteLesson).Methods(http.MethodDelete)

	// Student enrollment routes
	students := api.PathPrefix("/students").Subrouter()
	students.Use(middleware.AuthMiddleware(tokenProvider))
	students.HandleFunc("/{id}/courses", enrollmentHandler.GetStudentCourses).Methods(http.MethodGet)

	// Admin routes
	adminRoutes := api.NewRoute().Subrouter()
	adminRoutes.Use(middleware.AuthMiddleware(tokenProvider))
	adminRoutes.Use(middleware.RequireRole(entity.UserRoleAdmin))

	return r
}
