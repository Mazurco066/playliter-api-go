package account

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID           uint    `gorm:"primaryKey"`
	Avatar       *string `gorm:"default:'https://res.cloudinary.com/r4kta/image/upload/v1653796384/playliter/avatar/sample_capr2m.jpg'"`
	Email        string  `gorm:"unique"`
	IsEmailValid bool    `gorm:"default:false"`
	Username     string
	Name         string
	Role         string `gorm:"default:'player'"`
	Password     string
	IsActive     bool `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
