package bandusecase

import (
	bandrepo "github.com/mazurco066/playliter-api-go/data/repositories/band"
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
)

type BandRequestUseCase interface {
	Create(*band.BandRequest) error
	InviteExists(*account.Account, *band.Band) bool
	FindByAccount(*account.Account, *commoninputs.PagingParams) ([]*band.BandRequest, error)
	FindById(uint) (*band.BandRequest, error)
	Update(*band.BandRequest) error
}

type bandRequestUseCase struct {
	Repo bandrepo.BandRequestRepo
}

func NewBandRequestUseCase(repo bandrepo.BandRequestRepo) BandRequestUseCase {
	return &bandRequestUseCase{
		Repo: repo,
	}
}

func (uc *bandRequestUseCase) Create(request *band.BandRequest) error {
	return uc.Repo.Create(request)
}

func (uc *bandRequestUseCase) FindByAccount(a *account.Account, p *commoninputs.PagingParams) ([]*band.BandRequest, error) {
	if p.Limit == 0 {
		p.Limit = 100
	}
	return uc.Repo.FindByAccount(a, p)
}

func (uc *bandRequestUseCase) FindById(id uint) (*band.BandRequest, error) {
	return uc.Repo.FindById(id)
}

func (uc *bandRequestUseCase) InviteExists(a *account.Account, b *band.Band) bool {
	invite, _ := uc.Repo.FindByAccountAndBand(a, b)
	if invite != nil {
		return true
	}
	return false
}

func (uc *bandRequestUseCase) Update(request *band.BandRequest) error {
	return uc.Repo.Update(request)
}
