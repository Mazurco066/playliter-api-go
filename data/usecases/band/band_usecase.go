package bandusecase

import (
	bandrepo "github.com/mazurco066/playliter-api-go/data/repositories/band"
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
)

type BandUseCase interface {
	Create(*band.Band) error
	CreateInvite(*band.BandRequest) error
	CreateMember(*band.Member) error
	InviteExists(*account.Account, *band.Band) bool
	FindInviteById(uint) (*band.BandRequest, error)
	FindByAccount(*account.Account, *commoninputs.PagingParams) ([]*band.Band, error)
	FindById(uint) (*band.Band, error)
	FindMemberById(uint) (*band.Member, error)
	RemoveMember(*band.Member) error
	UpdateInvite(*band.BandRequest) error
	UpdateMember(*band.Member) error
}

type bandUseCase struct {
	Repo            bandrepo.Repo
	BandRequestRepo bandrepo.BandRequestRepo
	MemberRepo      bandrepo.MemberRepo
}

func NewBandUseCase(
	repo bandrepo.Repo,
	bandRequestRepo bandrepo.BandRequestRepo,
	memberRepo bandrepo.MemberRepo,
) BandUseCase {
	return &bandUseCase{
		Repo:            repo,
		BandRequestRepo: bandRequestRepo,
		MemberRepo:      memberRepo,
	}
}

func (uc *bandUseCase) Create(band *band.Band) error {
	return uc.Repo.Create(band)
}

func (uc *bandUseCase) CreateInvite(request *band.BandRequest) error {
	return uc.BandRequestRepo.Create(request)
}

func (uc *bandUseCase) CreateMember(member *band.Member) error {
	return uc.MemberRepo.Create(member)
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

func (uc *bandUseCase) FindInviteById(id uint) (*band.BandRequest, error) {
	return uc.BandRequestRepo.FindById(id)
}

func (uc *bandUseCase) FindMemberById(id uint) (*band.Member, error) {
	result, err := uc.MemberRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *bandUseCase) InviteExists(a *account.Account, b *band.Band) bool {
	invite, _ := uc.BandRequestRepo.FindByAccountAndBand(a, b)
	if invite != nil {
		return true
	}
	return false
}

func (uc *bandUseCase) RemoveMember(m *band.Member) error {
	return uc.MemberRepo.Remove(m)
}

func (uc *bandUseCase) UpdateInvite(request *band.BandRequest) error {
	return uc.BandRequestRepo.Update(request)
}

func (uc *bandUseCase) UpdateMember(m *band.Member) error {
	return uc.MemberRepo.Update(m)
}
