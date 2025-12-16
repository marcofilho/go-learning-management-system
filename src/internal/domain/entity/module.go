package entity

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CourseID   string         `gorm:"type:uuid;not null;index" json:"course_id"`
	Title      string         `gorm:"not null" json:"title"`
	OrderIndex int            `gorm:"not null;default:0" json:"order_index"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Course  Course          `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
	Lessons []LessonVersion `gorm:"foreignKey:ModuleID" json:"lessons,omitempty"`
}

func (Module) TableName() string {
	return "modules"
}
