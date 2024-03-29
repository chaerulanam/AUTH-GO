package migration

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Email          string `gorm:"size:255;index:email,unique"`
		Username       string `gorm:"size:255;index:username,unique"`
		Phone          string `gorm:"size:15;index:phone,unique"`
		PasswordHash   string `gorm:"size:255"`
		ResetAt        time.Time
		ResetExpire    time.Time
		ActivateHash   *string     `gorm:"size:255"`
		Status         *string     `gorm:"size:255"`
		StatusMessage  *string     `gorm:"size:255"`
		Active         bool        `gorm:"default:0"`
		ForcePassReset bool        `gorm:"default:0"`
		AuthLogin      []AuthLogin `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	AuthLogin struct {
		// gorm.Model
		ID        uint
		IPAddress string `gorm:"size:255"`
		UserID    uint
		Date      time.Time
		Success   int `gorm:"default:0"`
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
		ID          uint
		Name        string `gorm:"size:255"`
		Description string `gorm:"size:255"`
	}

	AuthPermission struct {
		// gorm.Model
		ID          uint
		Name        string `gorm:"size:255"`
		Description string `gorm:"size:255"`
	}

	AuthGroupPermission struct {
		// gorm.Model
		ID             uint
		GroupID        uint
		PermissionID   uint
		AuthGroup      AuthGroup      `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		AuthPermission AuthPermission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	AuthGroupUser struct {
		// gorm.Model
		ID        uint
		GroupID   uint
		UserID    uint
		User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		AuthGroup AuthGroup `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	AuthUserPermission struct {
		// gorm.Model
		ID             uint
		PermissionID   uint
		UserID         uint
		User           User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		AuthPermission AuthPermission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	ApiKey struct {
		// gorm.Model
		ID     uint
		Token  string `gorm:"size:255"`
		Domain string `gorm:"size:255"`
		Expire time.Time
	}

	Tabler interface {
		TableName() string
	}
)

// func (User) TableName() string {
// 	return "auth_users"
// }
