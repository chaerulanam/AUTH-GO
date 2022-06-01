package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Email          string `gorm:"size:255;index:email,unique"`
		Username       string `gorm:"size:255;index:username,unique"`
		PasswordHash   string `gorm:"size:255"`
		ResetAt        sql.NullTime
		ResetExpire    sql.NullTime
		ActivateHash   sql.NullString `gorm:"size:255"`
		Status         sql.NullString `gorm:"size:255"`
		StatusMessage  sql.NullString `gorm:"size:255"`
		Active         bool           `gorm:"default:0"`
		ForcePassReset bool           `gorm:"default:0"`
	}

	AuthLogin struct {
		// gorm.Model
		ID        uint
		IPAddress string `gorm:"size:255"`
		UserID    uint
		Date      time.Time
		Success   int  `gorm:"default:0"`
		User      User `gorm:"foreignKey:UserID"`
	}

	AuthToken struct {
		// gorm.Model
		ID           uint
		Selector     string `gorm:"size:255"`
		HasValidator string `gorm:"size:255"`
		UserID       uint
		Expire       time.Time
		User         User `gorm:"foreignKey:UserID"`
	}

	AuthResetAttemp struct {
		// gorm.Model
		ID        uint
		Email     string `gorm:"size:255"`
		IPAddress string `gorm:"size:255"`
		UserAgent string `gorm:"size:255"`
		Token     string `gorm:"size:255"`
		CreatedAt time.Time
	}

	AuthActivationAttemp struct {
		// gorm.Model
		ID        uint
		IPAddress string `gorm:"size:255"`
		UserAgent string `gorm:"size:255"`
		Token     string `gorm:"size:255"`
		CreatedAt time.Time
	}

	AuthGroup struct {
		// gorm.Model
		ID          uint32
		Name        string `gorm:"size:255"`
		Description string `gorm:"size:255"`
	}

	AuthPermission struct {
		// gorm.Model
		ID          uint32
		Name        string `gorm:"size:255"`
		Description string `gorm:"size:255"`
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

func (AuthLogin) TableName() string {
	return "auth_users"
}
