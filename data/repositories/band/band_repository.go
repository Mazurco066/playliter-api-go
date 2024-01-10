package bandrepo

import (
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*band.Band) error
}

type BandRepo struct {
	db *gorm.DB
}

func NewBandRepo(db *gorm.DB) Repo {
	return &BandRepo{
		db: db,
	}
}

func (repo *BandRepo) Create(band *band.Band) error {
	return repo.db.Create(band).Error
}
