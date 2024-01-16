package bandrepo

import (
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*band.Band) error
	FindByAccount(*account.Account, *commoninputs.PagingParams) ([]*band.Band, error)
	FindById(uint) (*band.Band, error)
	Remove(*band.Band) error
	Update(*band.Band) error
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

func (repo *BandRepo) FindByAccount(a *account.Account, p *commoninputs.PagingParams) ([]*band.Band, error) {
	var results []*band.Band
	if err := repo.db.
		Where("owner_id = ? OR EXISTS (SELECT 1 FROM members WHERE bands.id = members.band_id AND members.account_id = ?)", a.ID, a.ID).
		Preload("Owner").
		Preload("Members").
		Limit(p.Limit).
		Offset(p.Offset).
		Group("bands.id").
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *BandRepo) FindById(id uint) (*band.Band, error) {
	var band band.Band
	if err := repo.db.
		Where("id = ?", id).
		Preload("Owner").
		Preload("Members").
		First(&band).Error; err != nil {
		return nil, err
	}
	return &band, nil
}

func (repo *BandRepo) Update(band *band.Band) error {
	return repo.db.Save(band).Error
}

func (repo *BandRepo) Remove(band *band.Band) error {
	return repo.db.Delete(band).Error
}
