package dto

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginationRequest struct {
	Limit  int `json:"limit" form:"limit" validate:"omitempty,gte=1,lte=100"`
	Offset int `json:"offset" form:"offset" validate:"omitempty,gte=0"`
}

func NewErrorResponse(err error, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:   err.Error(),
		Message: message,
	}
}
