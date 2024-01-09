package concert

import (
	"gorm.io/gorm"

	"github.com/mazurco066/playliter-api-go/domain/models/song"
)

type ConcertSong struct {
	gorm.Model
	ConcertID uint      `json:"concert_id"`
	Concert   Concert   `gorm:"foreignKey:ConcertID" json:"concert"`
	SongID    uint      `json:"song_id"`
	Song      song.Song `gorm:"foreignKey:SongID" json:"song"`
}
