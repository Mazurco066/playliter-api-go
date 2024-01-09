package auth

import (
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	Token              string          `json:"token"`
	ResetPasswordToken *string         `json:"reset_password_token"`
	AccountID          uint            `json:"account_id"`
	Account            account.Account `gorm:"foreignKey:AccountID" json:"account"`
}
