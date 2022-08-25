package migration

import (
	"gorm.io/gorm"
)

type (
	ProfilSekolah struct {
		gorm.Model
		NamaSekolah      string `gorm:"size:255"`
		AlamatSekolah    string `gorm:"size:255"`
		NoHpSekolah      string `gorm:"size:15"`
		IconSekolah      string `gorm:"type:text"`
		LogoSekolah      string `gorm:"type:text"`
		DeskripsiSekolah string `gorm:"type:text"`
		DomainSekolah    string `gorm:"size:255"`
		Active           bool
	}

	MasterFiturSekolah struct {
		gorm.Model
		NamaFitur string `gorm:"size:255"`
		Deskripsi string `gorm:"size:255"`
	}

	DataFiturSekolah struct {
		gorm.Model
		ProfilSekolahID      uint
		MasterFiturSekolahID uint
		ProfilSekolah        ProfilSekolah      `gorm:"foreignKey:ProfilSekolahID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		MasterFiturSekolah   MasterFiturSekolah `gorm:"foreignKey:MasterFiturSekolahID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	AksesGrupFiturSekolah struct {
		gorm.Model
		AuthGroupID     uint
		ProfilSekolahID uint
		AuthGroup       AuthGroup     `gorm:"foreignKey:AuthGroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		ProfilSekolah   ProfilSekolah `gorm:"foreignKey:ProfilSekolahID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		View            bool
		Insert          bool
		Update          bool
		Delete          bool
	}
)

func MigrateAll(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&AuthLogin{})
	db.AutoMigrate(&AuthResetAttemp{})
	db.AutoMigrate(&AuthActivationAttemp{})
	db.AutoMigrate(&AuthGroup{})
	db.AutoMigrate(&AuthPermission{})
	db.AutoMigrate(&AuthGroupPermission{})
	db.AutoMigrate(&AuthGroupUser{})
	db.AutoMigrate(&AuthUserPermission{})
	db.AutoMigrate(&ApiKey{})
	db.AutoMigrate(&ProfilSekolah{})
	db.AutoMigrate(&MasterFiturSekolah{})
	db.AutoMigrate(&DataFiturSekolah{})
	db.AutoMigrate(&AksesGrupFiturSekolah{})
}
