package account

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Avatar       *string `gorm:"default:'https://res.cloudinary.com/r4kta/image/upload/v1653796384/playliter/avatar/sample_capr2m.jpg'" json:"avatar"`
	Email        string  `gorm:"unique" json:"email"`
	IsEmailValid bool    `gorm:"default:false" json:"is_email_valid"`
	Username     string  `gorm:"unique" json:"username"`
	Name         string  `json:"name"`
	Role         string  `gorm:"default:'player'" json:"role"`
	Password     string  `json:"password"`
	IsActive     bool    `gorm:"default:true" json:"is_active"`
}
