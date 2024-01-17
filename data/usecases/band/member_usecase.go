package bandusecase

import (
	bandrepo "github.com/mazurco066/playliter-api-go/data/repositories/band"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
)

type MemberUseCase interface {
	Create(*band.Member) error
	FindById(uint) (*band.Member, error)
	Remove(*band.Member) error
	Update(*band.Member) error
}

type memberUseCase struct {
	Repo bandrepo.MemberRepo
}

func NewMemberUseCase(repo bandrepo.MemberRepo) MemberUseCase {
	return &memberUseCase{
		Repo: repo,
	}
}

func (uc *memberUseCase) Create(member *band.Member) error {
	return uc.Repo.Create(member)
}

func (uc *memberUseCase) FindById(id uint) (*band.Member, error) {
	result, err := uc.Repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *memberUseCase) Remove(m *band.Member) error {
	return uc.Repo.Remove(m)
}

func (uc *memberUseCase) Update(m *band.Member) error {
	return uc.Repo.Update(m)
}
