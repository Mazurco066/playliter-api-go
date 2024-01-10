package concertrepo

import "gorm.io/gorm"

type Repo interface {
}

type ConcertRepo struct {
	db *gorm.DB
}

func NewConcertRepo(db *gorm.DB) Repo {
	return &ConcertRepo{
		db: db,
	}
}
