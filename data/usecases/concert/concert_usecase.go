package concertusecase

import concertrepo "github.com/mazurco066/playliter-api-go/data/repositories/concert"

type ConcertUseCase interface {
}

type concertUseCase struct {
	Repo concertrepo.Repo
}

func NewConcertUseCase(
	repo concertrepo.Repo,
) ConcertUseCase {
	return &concertUseCase{
		Repo: repo,
	}
}
