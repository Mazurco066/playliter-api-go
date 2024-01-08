package song

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/band"
)

type Song struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Title       string
	Writter     string
	Tone        string
	Body        string
	EmbeddedUrl *string
	Category    *string
	BandID      uint
	Band        band.Band `gorm:"foreignKey:BandID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
