package band

import (
	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/account"
)

type Band struct {
	gorm.Model
	Logo        *string `gorm:"default:'https://res.cloudinary.com/r4kta/image/upload/v1663515679/playliter/logo/default_band_mklz55.png'"`
	Title       string
	Description string
	OwnerID     uint
	Owner       account.Account `gorm:"foreignKey:OwnerID"`
	Members     []Member
}
