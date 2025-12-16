package dto

import "time"

type EnrollmentDTO struct {
	StudentID      string    `json:"student_id"`
	CourseID       string    `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type EnrollCourseRequest struct {
	StudentID string `json:"student_id" validate:"omitempty,uuid"`
}

type EnrollmentFilterRequest struct {
	Status   string `json:"status" validate:"omitempty,oneof=active dropped completed"`
	DateFrom string `json:"date_from" validate:"omitempty"`
	DateTo   string `json:"date_to" validate:"omitempty"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}
