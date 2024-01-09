package band

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/account"
)

type Member struct {
	gorm.Model
	BandID    uint
	Band      Band `gorm:"foreignKey:BandID"`
	AccountID uint
	Account   account.Account `gorm:"foreignKey:AccountID"`
	Role      string
	JoinedAt  time.Time
}
