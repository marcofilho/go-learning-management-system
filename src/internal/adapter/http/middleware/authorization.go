package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
)

// RequireCourseOwnership ensures user is admin or the course instructor
func RequireCourseOwnership(courseRepo repository.CourseRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "User not authenticated")
				return
			}

			// Admin can access everything
			if claims.Role == entity.UserRoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			// Get course ID from URL
			vars := mux.Vars(r)
			courseIDStr := vars["id"]
			if courseIDStr == "" {
				courseIDStr = vars["courseId"]
			}

			if courseIDStr == "" {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "Course ID not found in request")
				return
			}

			courseID, err := uuid.Parse(courseIDStr)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "Invalid course ID format")
				return
			}

			course, err := courseRepo.GetByID(r.Context(), courseID)
			if err != nil {
				if err == entity.ErrNotFound {
					respondWithError(w, http.StatusNotFound, entity.ErrNotFound, "Course not found")
				} else {
					respondWithError(w, http.StatusInternalServerError, err, "Failed to verify course ownership")
				}
				return
			}

			if course.InstructorID != claims.UserID {
				respondWithError(w, http.StatusForbidden, entity.ErrInsufficientPermissions, "Only the course instructor can perform this action")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireEnrollment ensures user is enrolled in the course (for students) or owns the course (instructors)
func RequireEnrollment(courseRepo repository.CourseRepository, enrollmentRepo repository.EnrollmentRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "User not authenticated")
				return
			}

			if claims.Role == entity.UserRoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			vars := mux.Vars(r)
			courseIDStr := vars["id"]
			if courseIDStr == "" {
				courseIDStr = vars["courseId"]
			}

			if courseIDStr == "" {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "Course ID not found in request")
				return
			}

			courseID, err := uuid.Parse(courseIDStr)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "Invalid course ID format")
				return
			}

			course, err := courseRepo.GetByID(r.Context(), courseID)
			if err != nil {
				if err == entity.ErrNotFound {
					respondWithError(w, http.StatusNotFound, entity.ErrNotFound, "Course not found")
				} else {
					respondWithError(w, http.StatusInternalServerError, err, "Failed to fetch course")
				}
				return
			}

			if claims.Role == entity.UserRoleInstructor && course.InstructorID == claims.UserID {
				next.ServeHTTP(w, r)
				return
			}

			if claims.Role == entity.UserRoleStudent {
				enrollment, err := enrollmentRepo.GetByStudentAndCourse(r.Context(), claims.UserID, courseID)
				if err != nil || enrollment == nil || !enrollment.IsActive() {
					respondWithError(w, http.StatusForbidden, entity.ErrInsufficientPermissions, "You must be enrolled in this course to access its content")
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireModuleOwnership ensures user is admin or the course instructor (via module's course)
func RequireModuleOwnership(moduleRepo repository.ModuleRepository, courseRepo repository.CourseRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "User not authenticated")
				return
			}

			if claims.Role == entity.UserRoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			vars := mux.Vars(r)
			moduleIDStr := vars["id"]
			if moduleIDStr == "" {
				moduleIDStr = vars["moduleId"]
			}

			if moduleIDStr == "" {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "Module ID not found in request")
				return
			}

			moduleID, err := uuid.Parse(moduleIDStr)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "Invalid module ID format")
				return
			}

			module, err := moduleRepo.GetByID(r.Context(), moduleID)
			if err != nil {
				if err == entity.ErrNotFound {
					respondWithError(w, http.StatusNotFound, entity.ErrNotFound, "Module not found")
				} else {
					respondWithError(w, http.StatusInternalServerError, err, "Failed to verify module ownership")
				}
				return
			}

			course, err := courseRepo.GetByID(r.Context(), module.CourseID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err, "Failed to verify course ownership")
				return
			}

			if course.InstructorID != claims.UserID {
				respondWithError(w, http.StatusForbidden, entity.ErrInsufficientPermissions, "Only the course instructor can perform this action")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdminOrSelf ensures user is admin or accessing their own user resource
func RequireAdminOrSelf() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "User not authenticated")
				return
			}

			if claims.Role == entity.UserRoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			vars := mux.Vars(r)
			userIDStr := vars["id"]

			if userIDStr == "" {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "User ID not found in request")
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, entity.ErrInvalidInput, "Invalid user ID format")
				return
			}

			if claims.UserID != userID {
				respondWithError(w, http.StatusForbidden, entity.ErrInsufficientPermissions, "You can only access your own resources")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user claims from context
func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*auth.Claims)
	return claims, ok
}
