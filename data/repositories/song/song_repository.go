package songrepo

import "gorm.io/gorm"

type Repo interface {
}

type SongRepo struct {
	db *gorm.DB
}

func NewSongRepo(db *gorm.DB) Repo {
	return &SongRepo{
		db: db,
	}
}
