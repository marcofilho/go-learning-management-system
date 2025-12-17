package dto

import "time"

type CourseDTO struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	InstructorID    string    `json:"instructor_id"`
	DifficultyLevel string    `json:"difficulty_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateCourseRequest struct {
	Title           string `json:"title" validate:"required"`
	Description     string `json:"description"`
	DifficultyLevel string `json:"difficulty_level" validate:"omitempty,oneof=beginner intermediate advanced"`
	InstructorID    string `json:"instructor_id" validate:"omitempty,uuid"`
}

type UpdateCourseRequest struct {
	Title           string `json:"title" validate:"omitempty"`
	Description     string `json:"description" validate:"omitempty"`
	DifficultyLevel string `json:"difficulty_level" validate:"omitempty,oneof=beginner intermediate advanced"`
}

type CourseFilterRequest struct {
	InstructorID    string `json:"instructor_id"`
	DifficultyLevel string `json:"difficulty_level" validate:"omitempty,oneof=beginner intermediate advanced"`
	ActiveOnly      bool   `json:"active_only"`
	Limit           int    `json:"limit"`
	Offset          int    `json:"offset"`
}
