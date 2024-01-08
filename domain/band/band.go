package band

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/account"
)

type Band struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey"`
	Logo        *string `gorm:"default:'https://res.cloudinary.com/r4kta/image/upload/v1663515679/playliter/logo/default_band_mklz55.png'"`
	Title       string
	Description string
	OwnerID     uint
	Owner       account.Account `gorm:"foreignKey:OwnerID"`
	Members     []Member
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
