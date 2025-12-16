package dto

import "time"

type LessonVersionDTO struct {
	ID            string    `json:"id"`
	ModuleID      string    `json:"module_id"`
	VersionNumber int       `json:"version_number"`
	Content       string    `json:"content"`
	VideoURL      string    `json:"video_url,omitempty"`
	AttachmentURL string    `json:"attachment_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateLessonRequest struct {
	Content       string `json:"content" validate:"required"`
	VideoURL      string `json:"video_url" validate:"omitempty,url"`
	AttachmentURL string `json:"attachment_url" validate:"omitempty,url"`
}

type CreateLessonVersionRequest struct {
	Content       string `json:"content" validate:"required"`
	VideoURL      string `json:"video_url" validate:"omitempty,url"`
	AttachmentURL string `json:"attachment_url" validate:"omitempty,url"`
}
