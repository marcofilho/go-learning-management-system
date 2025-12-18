package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
)

type CertificationWebhookHandler struct {
	webhookUseCase *usecase.CertificationWebhookUseCase
	webhookSecret  string
}

func NewCertificationWebhookHandler(webhookUseCase *usecase.CertificationWebhookUseCase) *CertificationWebhookHandler {
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		secret = "default-webhook-secret-change-in-production"
	}

	return &CertificationWebhookHandler{
		webhookUseCase: webhookUseCase,
		webhookSecret:  secret,
	}
}

// ProcessCertificationWebhook godoc
// @Summary Process certification webhook
// @Description Receives certification results from external provider
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Signature header string true "HMAC signature"
// @Param payload body usecase.CertificationWebhookPayload true "Certification data"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/certification-webhook [post]
func (h *CertificationWebhookHandler) ProcessCertificationWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Validate HMAC signature
	signature := r.Header.Get("X-Webhook-Signature")
	if signature == "" {
		RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Missing webhook signature"})
		return
	}

	if !h.validateSignature(body, signature) {
		RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid webhook signature"})
		return
	}

	// Parse payload
	var payload usecase.CertificationWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	// Validate payload fields
	if payload.StudentID == "" || payload.CourseID == "" || payload.CertificationStatus == "" {
		RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
		return
	}

	if payload.CertificationStatus != "passed" && payload.CertificationStatus != "failed" {
		RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "certification_status must be 'passed' or 'failed'"})
		return
	}

	if payload.Score < 0 || payload.Score > 100 {
		RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "score must be between 0 and 100"})
		return
	}

	// Process certification
	if err := h.webhookUseCase.ProcessCertification(r.Context(), payload); err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Certification webhook processed successfully"})
}

func (h *CertificationWebhookHandler) validateSignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
