package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Email          string
		Username       string
		PasswordHash   string
		ResetAt        sql.NullTime
		ResetExpire    sql.NullTime
		ActivateHash   sql.NullString
		Status         sql.NullString
		StatusMessage  sql.NullString
		Active         bool
		ForcePassReset bool
	}

	AuthLogin struct {
		// gorm.Model
		ID        uint
		IPAddress string
		UserID    uint
		CreatedAt time.Time `gorm:"column:Date"`
		Success   int
		User      User `gorm:"foreignKey:UserID"`
	}

	AuthToken struct {
		// gorm.Model
		ID           uint
		Selector     string
		HasValidator string
		UserID       uint
		Expire       sql.NullTime
		User         User `gorm:"foreignKey:UserID"`
	}

	AuthResetAttemp struct {
		// gorm.Model
		ID        uint
		Email     string
		IPAddress string
		UserAgent string
		Token     string
		CreatedAt time.Time
	}

	AuthActivationAttemp struct {
		// gorm.Model
		ID        uint
		IPAddress string
		UserAgent string
		Token     string
		CreatedAt time.Time
	}

	AuthGroup struct {
		// gorm.Model
		ID          uint32
		Name        string
		Description string
	}

	AuthPermission struct {
		// gorm.Model
		ID          uint32
		Name        string
		Description string
	}

	AuthGroupPermission struct {
		// gorm.Model
		ID             uint32
		GroupID        uint32
		PermissionID   uint32
		AuthGroup      AuthGroup      `gorm:"foreignKey:GroupID"`
		AuthPermission AuthPermission `gorm:"foreignKey:PermissionID"`
	}

	AuthGroupUser struct {
		// gorm.Model
		ID        uint
		GroupID   uint32
		UserID    uint
		AuthGroup AuthGroup `gorm:"foreignKey:GroupID"`
		User      User      `gorm:"foreignKey:UserID"`
	}

	AuthUserPermission struct {
		// gorm.Model
		ID             uint
		PermissionID   uint32
		UserID         uint
		AuthPermission AuthPermission `gorm:"foreignKey:PermissionID"`
		User           User           `gorm:"foreignKey:UserID"`
	}

	Tabler interface {
		TableName() string
	}
)

func (User) TableName() string {
	return "auth_users"
}
