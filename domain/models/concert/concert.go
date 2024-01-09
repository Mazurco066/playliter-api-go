package concert

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"github.com/mazurco066/playliter-api-go/domain/models/song"
)

type Concert struct {
	gorm.Model
	Title       string
	Description string
	Date        time.Time
	BandID      uint
	Band        band.Band   `gorm:"foreignKey:BandID"`
	Songs       []song.Song `gorm:"many2many:concert_songs;"`
}
