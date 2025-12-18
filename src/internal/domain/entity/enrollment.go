package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentStatus string

const (
	EnrollmentStatusActive    EnrollmentStatus = "active"
	EnrollmentStatusDropped   EnrollmentStatus = "dropped"
	EnrollmentStatusCompleted EnrollmentStatus = "completed"
)

type Enrollment struct {
	StudentID      string           `gorm:"primaryKey;type:uuid;not null" json:"student_id"`
	CourseID       string           `gorm:"primaryKey;type:uuid;not null" json:"course_id"`
	EnrollmentDate time.Time        `gorm:"not null;autoCreateTime" json:"enrollment_date"`
	Status         EnrollmentStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`

	Student User   `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"-"`
	Course  Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Enrollment) TableName() string {
	return "course_enrollments"
}

func (e *Enrollment) Validate() error {
	if strings.TrimSpace(e.StudentID) == "" {
		return fmt.Errorf("%w: student_id is required", ErrFieldRequired)
	}

	if _, err := uuid.Parse(e.StudentID); err != nil {
		return fmt.Errorf("%w: student_id must be a valid UUID", ErrInvalidInput)
	}

	if strings.TrimSpace(e.CourseID) == "" {
		return fmt.Errorf("%w: course_id is required", ErrFieldRequired)
	}

	if _, err := uuid.Parse(e.CourseID); err != nil {
		return fmt.Errorf("%w: course_id must be a valid UUID", ErrInvalidInput)
	}

	if e.Status != EnrollmentStatusActive &&
		e.Status != EnrollmentStatusDropped &&
		e.Status != EnrollmentStatusCompleted {
		return fmt.Errorf("%w: status must be one of: active, dropped, completed", ErrInvalidInput)
	}

	return nil
}

func NewEnrollment(student *User, courseID string) (*Enrollment, error) {
	if !student.CanEnrollInCourses() {
		return nil, ErrInvalidInput
	}

	enrollment := &Enrollment{
		StudentID:      student.ID,
		CourseID:       courseID,
		EnrollmentDate: time.Now(),
		Status:         EnrollmentStatusActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := enrollment.Validate(); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (e *Enrollment) CanBeModifiedBy(user *User) bool {
	return user.IsAdmin() || e.StudentID == user.ID
}

func (e *Enrollment) IsActive() bool {
	return e.Status == EnrollmentStatusActive
}

func (e *Enrollment) IsCompleted() bool {
	return e.Status == EnrollmentStatusCompleted
}

func (e *Enrollment) Drop() {
	e.Status = EnrollmentStatusDropped
	e.UpdatedAt = time.Now()
}

func (e *Enrollment) Complete() {
	e.Status = EnrollmentStatusCompleted
	e.UpdatedAt = time.Now()
}
