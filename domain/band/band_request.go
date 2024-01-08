package band

import (
	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/account"
)

type BandRequest struct {
	gorm.Model
	BandID    uint
	Band      Band `gorm:"foreignKey:BandID"`
	InvitedID uint
	Invited   account.Account `gorm:"foreignKey:InvitedID"`
	Status    string          // "pending", "accepted", "denied"
}
