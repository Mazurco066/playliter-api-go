package bandusecase

import (
	bandrepo "github.com/mazurco066/playliter-api-go/data/repositories/band"
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
)

type BandUseCase interface {
	Create(*band.Band) error
	FindByAccount(*account.Account, *commoninputs.PagingParams) ([]*band.Band, error)
	FindById(uint) (*band.Band, error)
	Remove(*band.Band) error
	Update(*band.Band) error
}

type bandUseCase struct {
	Repo bandrepo.Repo
}

func NewBandUseCase(repo bandrepo.Repo) BandUseCase {
	return &bandUseCase{
		Repo: repo,
	}
}

func (uc *bandUseCase) Create(band *band.Band) error {
	return uc.Repo.Create(band)
}

func (uc *bandUseCase) FindByAccount(a *account.Account, p *commoninputs.PagingParams) ([]*band.Band, error) {
	if p.Limit == 0 {
		p.Limit = 100
	}
	return uc.Repo.FindByAccount(a, p)
}

func (uc *bandUseCase) FindById(id uint) (*band.Band, error) {
	result, err := uc.Repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *bandUseCase) Remove(b *band.Band) error {
	return uc.Repo.Remove(b)
}

func (uc *bandUseCase) Update(b *band.Band) error {
	return uc.Repo.Update(b)
}
