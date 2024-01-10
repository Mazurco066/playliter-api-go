package songusecase

import songrepo "github.com/mazurco066/playliter-api-go/data/repositories/song"

type SongUseCase interface {
}

type songUseCase struct {
	Repo songrepo.Repo
}

func NewSongUseCase(
	repo songrepo.Repo,
) SongUseCase {
	return &songUseCase{
		Repo: repo,
	}
}
