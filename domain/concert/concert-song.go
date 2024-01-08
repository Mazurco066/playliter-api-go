package concert

import (
	"time"

	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/song"
)

type ConcertSong struct {
	gorm.Model
	ConcertID uint
	Concert   Concert `gorm:"foreignKey:ConcertID"`
	SongID    uint
	Song      song.Song `gorm:"foreignKey:SongID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
