package concert

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/band"
	"github.com/mazurco066/playliter-api-go/domain/song"
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
