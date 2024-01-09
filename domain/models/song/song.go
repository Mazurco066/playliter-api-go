package song

import (
	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/band"
)

type Song struct {
	gorm.Model
	Title       string    `json:"title"`
	Writter     string    `json:"writter"`
	Tone        string    `json:"tone"`
	Body        string    `json:"body"`
	EmbeddedUrl *string   `json:"embedded_url"`
	Category    *string   `json:"category"`
	BandID      uint      `json:"band_id"`
	Band        band.Band `gorm:"foreignKey:BandID" json:"band"`
}
