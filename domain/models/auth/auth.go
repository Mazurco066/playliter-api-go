package auth

import (
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	Token              string
	ResetPasswordToken *string
	AccountID          uint
	Account            account.Account `gorm:"foreignKey:AccountID"`
}
