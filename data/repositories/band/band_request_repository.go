package bandrepo

import (
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"gorm.io/gorm"
)

type BandRequestRepo interface {
	Create(*band.BandRequest) error
	FindByAccountAndBand(*account.Account, *band.Band) (*band.BandRequest, error)
}

type bandRequestRepo struct {
	db *gorm.DB
}

func NewBandRequestRepo(db *gorm.DB) BandRequestRepo {
	return &bandRequestRepo{
		db: db,
	}
}

func (repo *bandRequestRepo) Create(request *band.BandRequest) error {
	return repo.db.Create(request).Error
}

func (repo *bandRequestRepo) FindByAccountAndBand(a *account.Account, b *band.Band) (*band.BandRequest, error) {
	var request band.BandRequest
	if err := repo.db.
		Where("band_id = ? AND invited_id = ? AND status = ?", b.ID, a.ID, "pending").
		Preload("Band").
		Preload("Invited").
		First(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}
