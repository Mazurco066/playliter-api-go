package bandusecase

import (
	bandrepo "github.com/mazurco066/playliter-api-go/data/repositories/band"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
)

type BandUseCase interface {
	Create(*band.Band) error
}

type bandUseCase struct {
	Repo bandrepo.Repo
}

func NewBandUseCase(
	repo bandrepo.Repo,
) BandUseCase {
	return &bandUseCase{
		Repo: repo,
	}
}

func (uc *bandUseCase) Create(band *band.Band) error {
	return uc.Repo.Create(band)
}
