package accountusecase

import (
	accountrepo "github.com/mazurco066/playliter-api-go/data/repositories/account"
)

type AccountUseCase interface {
}

type accountUseCase struct {
	Repo accountrepo.Repo
}

func NewAccountUseCase(repo accountrepo.Repo) AccountUseCase {
	return &accountUseCase{
		Repo: repo,
	}
}
