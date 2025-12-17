package main

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/docs"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/handler"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
)

func SetupRoutes(
	userHandler *handler.UserHandler,
	courseHandler *handler.CourseHandler,
	enrollmentHandler *handler.EnrollmentHandler,
	moduleHandler *handler.ModuleHandler,
	lessonHandler *handler.LessonHandler,
	tokenProvider auth.TokenProvider,
	courseRepo repository.CourseRepository,
	enrollmentRepo repository.EnrollmentRepository,
	moduleRepo repository.ModuleRepository,
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

	// Admin only - list all users
	users.Handle("", middleware.RequireRole(entity.UserRoleAdmin)(http.HandlerFunc(userHandler.ListUsers))).Methods(http.MethodGet)

	// Any authenticated user can view their own profile, admin can view all
	users.HandleFunc("/{id}", userHandler.GetUserByID).Methods(http.MethodGet)

	// Users can update their own profile, admin can update any
	users.Handle("/{id}", middleware.RequireAdminOrSelf()(http.HandlerFunc(userHandler.UpdateUser))).Methods(http.MethodPut)

	// Admin only - delete users
	users.Handle("/{id}", middleware.RequireRole(entity.UserRoleAdmin)(http.HandlerFunc(userHandler.DeleteUser))).Methods(http.MethodDelete)

	// Course routes
	courses := api.PathPrefix("/courses").Subrouter()
	courses.Use(middleware.AuthMiddleware(tokenProvider))

	// All authenticated users can list and view courses
	courses.HandleFunc("", courseHandler.ListCourses).Methods(http.MethodGet)
	courses.HandleFunc("/{id}", courseHandler.GetCourse).Methods(http.MethodGet)

	// Admin and instructors can create courses
	courses.Handle("", middleware.RequireRole(entity.UserRoleAdmin, entity.UserRoleInstructor)(http.HandlerFunc(courseHandler.CreateCourse))).Methods(http.MethodPost)

	// Admin or course owner can update/delete courses
	courses.Handle("/{id}", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(courseHandler.UpdateCourse))).Methods(http.MethodPut)
	courses.Handle("/{id}", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(courseHandler.DeleteCourse))).Methods(http.MethodDelete)

	// Course enrollment routes
	// Students can self-enroll, admin can enroll anyone
	courses.HandleFunc("/{id}/enroll", enrollmentHandler.EnrollInCourse).Methods(http.MethodPost)

	// Admin or course instructor can view enrolled students
	courses.Handle("/{id}/students", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(enrollmentHandler.GetCourseStudents))).Methods(http.MethodGet)

	// Course modules routes
	// Admin or course instructor can create modules
	courses.Handle("/{courseId}/modules", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(moduleHandler.CreateModule))).Methods(http.MethodPost)

	// Enrolled students, course instructor, or admin can view modules
	courses.Handle("/{courseId}/modules", middleware.RequireEnrollment(courseRepo, enrollmentRepo)(http.HandlerFunc(moduleHandler.GetCourseModules))).Methods(http.MethodGet)

	// Module routes
	modules := api.PathPrefix("/modules").Subrouter()
	modules.Use(middleware.AuthMiddleware(tokenProvider))

	// Admin or course instructor can update/delete modules
	modules.Handle("/{id}", middleware.RequireModuleOwnership(moduleRepo, courseRepo)(http.HandlerFunc(moduleHandler.UpdateModule))).Methods(http.MethodPut)
	modules.Handle("/{id}", middleware.RequireModuleOwnership(moduleRepo, courseRepo)(http.HandlerFunc(moduleHandler.DeleteModule))).Methods(http.MethodDelete)

	// Module lesson routes
	// Admin or course instructor can create lessons
	modules.Handle("/{moduleId}/lessons", middleware.RequireModuleOwnership(moduleRepo, courseRepo)(http.HandlerFunc(lessonHandler.CreateLesson))).Methods(http.MethodPost)

	// Enrolled students, course instructor, or admin can view lessons (will be checked in handler)
	modules.HandleFunc("/{moduleId}/lessons", lessonHandler.GetModuleLessons).Methods(http.MethodGet)

	// Lesson routes
	lessons := api.PathPrefix("/lessons").Subrouter()
	lessons.Use(middleware.AuthMiddleware(tokenProvider))

	// Admin or course instructor can create lesson versions and delete lessons
	// Note: Ownership check will be done in handler based on lesson's course
	lessons.Handle("/{lessonId}/version", middleware.RequireRole(entity.UserRoleAdmin, entity.UserRoleInstructor)(http.HandlerFunc(lessonHandler.CreateLessonVersion))).Methods(http.MethodPost)
	lessons.Handle("/{lessonId}", middleware.RequireRole(entity.UserRoleAdmin, entity.UserRoleInstructor)(http.HandlerFunc(lessonHandler.DeleteLesson))).Methods(http.MethodDelete)

	// Admin only - view all lesson versions (audit trail)
	lessons.Handle("/{lessonId}/all-versions", middleware.RequireRole(entity.UserRoleAdmin)(http.HandlerFunc(lessonHandler.GetAllLessonVersions))).Methods(http.MethodGet)

	// Student enrollment routes
	students := api.PathPrefix("/students").Subrouter()
	students.Use(middleware.AuthMiddleware(tokenProvider))

	// Users can view their own enrollments, admin can view all
	students.Handle("/{id}/courses", middleware.RequireAdminOrSelf()(http.HandlerFunc(enrollmentHandler.GetStudentCourses))).Methods(http.MethodGet)

	return r
}
