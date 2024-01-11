package band

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/account"
)

type Member struct {
	gorm.Model
	BandID    uint            `json:"band_id"`
	Band      Band            `gorm:"foreignKey:BandID" json:"band"`
	AccountID uint            `json:"account_id"`
	Account   account.Account `gorm:"foreignKey:AccountID" json:"account"`
	Role      string          `json:"role"` // "member", "admin"
	JoinedAt  time.Time       `json:"joined_at"`
}
