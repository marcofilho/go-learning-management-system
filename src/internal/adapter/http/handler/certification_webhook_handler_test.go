package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func generateValidSignature(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_Success(t *testing.T) {
	secret := "test-secret"
	os.Setenv("WEBHOOK_SECRET", secret)
	defer os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	studentID := uuid.New()
	courseID := uuid.New()
	payload := usecase.CertificationWebhookPayload{
		StudentID:           studentID,
		CourseID:            courseID,
		CertificationStatus: "passed",
		Score:               95,
	}

	body, _ := json.Marshal(payload)
	signature := generateValidSignature(body, secret)

	student := &entity.User{ID: studentID}
	course := &entity.Course{ID: courseID}
	enrollment := &entity.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
		Status:    entity.EnrollmentStatusActive,
	}

	req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
	req.Header.Set("X-Webhook-Signature", signature)
	rec := httptest.NewRecorder()

	mockUserRepo.On("GetByID", mock.Anything, studentID).Return(student, nil)
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil)
	mockEnrollmentRepo.On("GetByStudentAndCourse", mock.Anything, studentID, courseID).Return(enrollment, nil)
	mockEnrollmentRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
	mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	handler.ProcessCertificationWebhook(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUserRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
	mockAuditLogRepo.AssertExpectations(t)
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_MissingSignature(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	payload := usecase.CertificationWebhookPayload{
		StudentID:           uuid.New(),
		CourseID:            uuid.New(),
		CertificationStatus: "passed",
		Score:               95,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ProcessCertificationWebhook(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Missing webhook signature")
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_InvalidSignature(t *testing.T) {
	secret := "test-secret"
	os.Setenv("WEBHOOK_SECRET", secret)
	defer os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	payload := usecase.CertificationWebhookPayload{
		StudentID:           uuid.New(),
		CourseID:            uuid.New(),
		CertificationStatus: "passed",
		Score:               95,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
	req.Header.Set("X-Webhook-Signature", "invalid-signature")
	rec := httptest.NewRecorder()

	handler.ProcessCertificationWebhook(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid webhook signature")
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_InvalidJSON(t *testing.T) {
	secret := "test-secret"
	os.Setenv("WEBHOOK_SECRET", secret)
	defer os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	body := []byte("invalid json")
	signature := generateValidSignature(body, secret)

	req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
	req.Header.Set("X-Webhook-Signature", signature)
	rec := httptest.NewRecorder()

	handler.ProcessCertificationWebhook(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid JSON payload")
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_MissingFields(t *testing.T) {
	secret := "test-secret"
	os.Setenv("WEBHOOK_SECRET", secret)
	defer os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	tests := []struct {
		name    string
		payload usecase.CertificationWebhookPayload
	}{
		{
			name: "missing student ID",
			payload: usecase.CertificationWebhookPayload{
				CourseID:            uuid.New(),
				CertificationStatus: "passed",
				Score:               95,
			},
		},
		{
			name: "missing course ID",
			payload: usecase.CertificationWebhookPayload{
				StudentID:           uuid.New(),
				CertificationStatus: "passed",
				Score:               95,
			},
		},
		{
			name: "missing certification status",
			payload: usecase.CertificationWebhookPayload{
				StudentID: uuid.New(),
				CourseID:  uuid.New(),
				Score:     95,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			signature := generateValidSignature(body, secret)

			req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
			req.Header.Set("X-Webhook-Signature", signature)
			rec := httptest.NewRecorder()

			handler.ProcessCertificationWebhook(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Missing required fields")
		})
	}
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_InvalidStatus(t *testing.T) {
	secret := "test-secret"
	os.Setenv("WEBHOOK_SECRET", secret)
	defer os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	payload := usecase.CertificationWebhookPayload{
		StudentID:           uuid.New(),
		CourseID:            uuid.New(),
		CertificationStatus: "invalid",
		Score:               95,
	}

	body, _ := json.Marshal(payload)
	signature := generateValidSignature(body, secret)

	req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
	req.Header.Set("X-Webhook-Signature", signature)
	rec := httptest.NewRecorder()

	handler.ProcessCertificationWebhook(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "certification_status must be")
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_InvalidScore(t *testing.T) {
	secret := "test-secret"
	os.Setenv("WEBHOOK_SECRET", secret)
	defer os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	tests := []struct {
		name  string
		score int
	}{
		{"score too low", -1},
		{"score too high", 101},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := usecase.CertificationWebhookPayload{
				StudentID:           uuid.New(),
				CourseID:            uuid.New(),
				CertificationStatus: "passed",
				Score:               tt.score,
			}

			body, _ := json.Marshal(payload)
			signature := generateValidSignature(body, secret)

			req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
			req.Header.Set("X-Webhook-Signature", signature)
			rec := httptest.NewRecorder()

			handler.ProcessCertificationWebhook(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "score must be between 0 and 100")
		})
	}
}

func TestCertificationWebhookHandler_ProcessCertificationWebhook_UseCaseError(t *testing.T) {
	secret := "test-secret"
	os.Setenv("WEBHOOK_SECRET", secret)
	defer os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	studentID := uuid.New()
	courseID := uuid.New()
	payload := usecase.CertificationWebhookPayload{
		StudentID:           studentID,
		CourseID:            courseID,
		CertificationStatus: "passed",
		Score:               95,
	}

	body, _ := json.Marshal(payload)
	signature := generateValidSignature(body, secret)

	req := httptest.NewRequest(http.MethodPost, "/api/certification-webhook", bytes.NewReader(body))
	req.Header.Set("X-Webhook-Signature", signature)
	rec := httptest.NewRecorder()

	// Student not found
	mockUserRepo.On("GetByID", mock.Anything, studentID).Return(nil, entity.ErrNotFound)

	handler.ProcessCertificationWebhook(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestNewCertificationWebhookHandler_DefaultSecret(t *testing.T) {
	os.Unsetenv("WEBHOOK_SECRET")

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditLogRepo := new(MockAuditLogRepository)

	webhookUC := usecase.NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditLogRepo)
	handler := NewCertificationWebhookHandler(webhookUC)

	assert.Equal(t, "default-webhook-secret-change-in-production", handler.webhookSecret)
}
