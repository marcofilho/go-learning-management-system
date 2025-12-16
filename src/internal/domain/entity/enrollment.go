package entity

import (
	"time"

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

func (e *Enrollment) IsActive() bool {
	return e.Status == EnrollmentStatusActive
}

func (e *Enrollment) IsCompleted() bool {
	return e.Status == EnrollmentStatusCompleted
}

func (e *Enrollment) Drop() {
	e.Status = EnrollmentStatusDropped
}

func (e *Enrollment) Complete() {
	e.Status = EnrollmentStatusCompleted
}
