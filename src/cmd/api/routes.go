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
	auditLogHandler *handler.AuditLogHandler,
	certificationWebhookHandler *handler.CertificationWebhookHandler,
	tokenProvider auth.TokenProvider,
	courseRepo repository.CourseRepository,
	enrollmentRepo repository.EnrollmentRepository,
	moduleRepo repository.ModuleRepository,
) http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	api.HandleFunc("/certification-webhook", certificationWebhookHandler.ProcessCertificationWebhook).Methods(http.MethodPost)

	api.Handle("/auth/register", middleware.OptionalAuthMiddleware(tokenProvider)(http.HandlerFunc(userHandler.Register))).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", userHandler.Login).Methods(http.MethodPost)

	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware(tokenProvider))
	users.Handle("", middleware.RequireRole(entity.UserRoleAdmin)(http.HandlerFunc(userHandler.ListUsers))).Methods(http.MethodGet)
	users.HandleFunc("/{id}", userHandler.GetUserByID).Methods(http.MethodGet)
	users.Handle("/{id}", middleware.RequireAdminOrSelf()(http.HandlerFunc(userHandler.UpdateUser))).Methods(http.MethodPut)
	users.Handle("/{id}", middleware.RequireRole(entity.UserRoleAdmin)(http.HandlerFunc(userHandler.DeleteUser))).Methods(http.MethodDelete)

	courses := api.PathPrefix("/courses").Subrouter()
	courses.Use(middleware.AuthMiddleware(tokenProvider))
	courses.HandleFunc("", courseHandler.ListCourses).Methods(http.MethodGet)
	courses.HandleFunc("/{id}", courseHandler.GetCourse).Methods(http.MethodGet)
	courses.Handle("", middleware.RequireRole(entity.UserRoleAdmin, entity.UserRoleInstructor)(http.HandlerFunc(courseHandler.CreateCourse))).Methods(http.MethodPost)
	courses.Handle("/{id}", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(courseHandler.UpdateCourse))).Methods(http.MethodPut)
	courses.Handle("/{id}", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(courseHandler.DeleteCourse))).Methods(http.MethodDelete)

	courses.HandleFunc("/{id}/enroll", enrollmentHandler.EnrollInCourse).Methods(http.MethodPost)
	courses.Handle("/{id}/students", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(enrollmentHandler.GetCourseStudents))).Methods(http.MethodGet)

	courses.Handle("/{courseId}/modules", middleware.RequireCourseOwnership(courseRepo)(http.HandlerFunc(moduleHandler.CreateModule))).Methods(http.MethodPost)
	courses.Handle("/{courseId}/modules", middleware.RequireEnrollment(courseRepo, enrollmentRepo)(http.HandlerFunc(moduleHandler.GetCourseModules))).Methods(http.MethodGet)

	modules := api.PathPrefix("/modules").Subrouter()
	modules.Use(middleware.AuthMiddleware(tokenProvider))
	modules.Handle("/{id}", middleware.RequireModuleOwnership(moduleRepo, courseRepo)(http.HandlerFunc(moduleHandler.UpdateModule))).Methods(http.MethodPut)
	modules.Handle("/{id}", middleware.RequireModuleOwnership(moduleRepo, courseRepo)(http.HandlerFunc(moduleHandler.DeleteModule))).Methods(http.MethodDelete)

	modules.Handle("/{moduleId}/lessons", middleware.RequireModuleOwnership(moduleRepo, courseRepo)(http.HandlerFunc(lessonHandler.CreateLesson))).Methods(http.MethodPost)
	modules.Handle("/{moduleId}/lessons", middleware.RequireModuleAccess(moduleRepo, courseRepo, enrollmentRepo)(http.HandlerFunc(lessonHandler.GetModuleLessons))).Methods(http.MethodGet)

	lessons := api.PathPrefix("/lessons").Subrouter()
	lessons.Use(middleware.AuthMiddleware(tokenProvider))
	lessons.Handle("/{lessonId}/version", middleware.RequireRole(entity.UserRoleAdmin, entity.UserRoleInstructor)(http.HandlerFunc(lessonHandler.CreateLessonVersion))).Methods(http.MethodPost)
	lessons.Handle("/{lessonId}", middleware.RequireRole(entity.UserRoleAdmin, entity.UserRoleInstructor)(http.HandlerFunc(lessonHandler.DeleteLesson))).Methods(http.MethodDelete)
	lessons.Handle("/{lessonId}/all-versions", middleware.RequireRole(entity.UserRoleAdmin)(http.HandlerFunc(lessonHandler.GetAllLessonVersions))).Methods(http.MethodGet)

	students := api.PathPrefix("/students").Subrouter()
	students.Use(middleware.AuthMiddleware(tokenProvider))
	students.Handle("/{id}/courses", middleware.RequireAdminOrSelf()(http.HandlerFunc(enrollmentHandler.GetStudentCourses))).Methods(http.MethodGet)

	auditLogs := api.PathPrefix("/audit-logs").Subrouter()
	auditLogs.Use(middleware.AuthMiddleware(tokenProvider))
	auditLogs.Handle("", middleware.RequireRole(entity.UserRoleAdmin)(http.HandlerFunc(auditLogHandler.ListAuditLogs))).Methods(http.MethodGet)

	return r
}
