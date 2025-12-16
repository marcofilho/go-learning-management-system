package dto

import "time"

type ModuleDTO struct {
	ID         string    `json:"id"`
	CourseID   string    `json:"course_id"`
	Title      string    `json:"title"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateModuleRequest struct {
	Title      string `json:"title" validate:"required"`
	OrderIndex int    `json:"order_index" validate:"gte=0"`
}

type UpdateModuleRequest struct {
	Title      string `json:"title" validate:"omitempty"`
	OrderIndex int    `json:"order_index" validate:"omitempty,gte=0"`
}
