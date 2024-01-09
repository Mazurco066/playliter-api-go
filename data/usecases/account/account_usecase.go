package accountusecase

import (
	accountrepo "github.com/mazurco066/playliter-api-go/data/repositories/account"
	"github.com/mazurco066/playliter-api-go/domain/account"
)

type AccountUseCase interface {
	Create(*account.Account) error
}

type accountUseCase struct {
	Repo accountrepo.Repo
}

func NewAccountUseCase(repo accountrepo.Repo) AccountUseCase {
	return &accountUseCase{
		Repo: repo,
	}
}

func (uc *accountUseCase) Create(account *account.Account) error {
	return uc.Repo.Create(account)
}
