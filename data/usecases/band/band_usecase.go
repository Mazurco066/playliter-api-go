package bandusecase

import bandrepo "github.com/mazurco066/playliter-api-go/data/repositories/band"

type BandUseCase interface {
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
