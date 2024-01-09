package band

import (
	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/account"
)

type BandRequest struct {
	gorm.Model
	BandID    uint            `json:"band_id"`
	Band      Band            `gorm:"foreignKey:BandID" json:"band"`
	InvitedID uint            `json:"invited_id"`
	Invited   account.Account `gorm:"foreignKey:InvitedID" json:"invited"`
	Status    string          `json:"status"` // "pending", "accepted", "denied"
}
