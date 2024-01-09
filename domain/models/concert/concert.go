package concert

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"github.com/mazurco066/playliter-api-go/domain/models/song"
)

type Concert struct {
	gorm.Model
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Date        time.Time   `json:"date"`
	BandID      uint        `json:"band_id"`
	Band        band.Band   `gorm:"foreignKey:BandID" json:"band"`
	Songs       []song.Song `gorm:"many2many:concert_songs;" json:"songs"`
}
