package entity

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleStudent    UserRole = "student"
	UserRoleInstructor UserRole = "instructor"
	UserRoleAdmin      UserRole = "admin"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type User struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FirstName string         `gorm:"not null" json:"first_name"`
	LastName  string         `gorm:"not null" json:"last_name"`
	Role      UserRole       `gorm:"type:varchar(20);not null;default:'student'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("%w: email is required", ErrFieldRequired)
	}

	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("%w: invalid email format", ErrInvalidInput)
	}

	if strings.TrimSpace(u.FirstName) == "" {
		return fmt.Errorf("%w: first_name is required", ErrFieldRequired)
	}

	if strings.TrimSpace(u.LastName) == "" {
		return fmt.Errorf("%w: last_name is required", ErrFieldRequired)
	}

	if u.Role != UserRoleStudent && u.Role != UserRoleInstructor && u.Role != UserRoleAdmin {
		return fmt.Errorf("%w: role must be one of: student, instructor, admin", ErrInvalidInput)
	}

	return nil
}

func NewUser(email, password, firstName, lastName string, role UserRole) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) CanAuthenticate() error {
	if !u.IsActive {
		return ErrUnauthorized
	}
	return nil
}

func (u *User) CanEnrollInCourses() bool {
	return u.Role == UserRoleStudent && u.IsActive
}

func (u *User) IsInstructor() bool {
	return u.Role == UserRoleInstructor || u.Role == UserRoleAdmin
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) IsOwnedBy(id uuid.UUID) bool {
	return u.ID == id
}

func (u *User) CanCreateCourses() bool {
	return u.Role == UserRoleInstructor || u.Role == UserRoleAdmin
}
