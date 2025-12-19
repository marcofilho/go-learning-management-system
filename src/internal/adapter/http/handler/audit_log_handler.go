package handler

import (
	"net/http"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
)

type AuditLogHandler struct {
	auditLogRepo repository.AuditLogRepository
}

func NewAuditLogHandler(auditLogRepo repository.AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{auditLogRepo: auditLogRepo}
}

// ListAuditLogs godoc
// @Summary List audit logs
// @Description List audit logs (admin only)
// @Tags AuditLogs
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /audit-logs [get]
func (h *AuditLogHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	pg := parsePagination(r)

	logs, total, err := h.auditLogRepo.List(r.Context(), pg.limit, pg.offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch audit logs")
		return
	}

	respondWithJSON(w, http.StatusOK, paginatedResponse{
		Data:       logs,
		Pagination: buildPagination(pg.page, pg.pageSize, total),
	})
}
