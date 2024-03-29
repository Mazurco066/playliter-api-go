package accountusecase

import (
	accountrepo "github.com/mazurco066/playliter-api-go/data/repositories/account"
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/infra/hmachash"
	"golang.org/x/crypto/bcrypt"
)

type AccountUseCase interface {
	ComparePassword(rawPassword string, passwordFromDB string) error
	Create(*account.Account) error
	GetAccountById(uint) (*account.Account, error)
	GetAccountByEmail(string) (*account.Account, error)
	GetAccountByUsernameOrEmail(string) (*account.Account, error)
	HashPassword(string) (string, error)
	ListActiveAccounts(*account.Account, *commoninputs.PagingParams) ([]*account.Account, error)
	Update(*account.Account, bool) error
}

type accountUseCase struct {
	Repo accountrepo.Repo
	hmac hmachash.HMAC
}

func NewAccountUseCase(
	repo accountrepo.Repo,
	hmac hmachash.HMAC,
) AccountUseCase {
	return &accountUseCase{
		Repo: repo,
		hmac: hmac,
	}
}

func (uc *accountUseCase) ComparePassword(rawPassword string, passwordFromDB string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(passwordFromDB),
		[]byte(rawPassword),
	)
}

func (uc *accountUseCase) Create(account *account.Account) error {
	hashedPass, err := uc.HashPassword(account.Password)
	if err != nil {
		return err
	}
	account.Password = hashedPass
	return uc.Repo.Create(account)
}

func (uc *accountUseCase) GetAccountById(id uint) (*account.Account, error) {
	result, err := uc.Repo.FindByUserID(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *accountUseCase) GetAccountByEmail(email string) (*account.Account, error) {
	result, err := uc.Repo.FindByUserEmail(email)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *accountUseCase) GetAccountByUsernameOrEmail(filter string) (*account.Account, error) {
	result, err := uc.Repo.FindByUsernameOrEmail(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *accountUseCase) HashPassword(rawPassword string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), err
}

func (uc *accountUseCase) ListActiveAccounts(account *account.Account, paging *commoninputs.PagingParams) ([]*account.Account, error) {
	if paging.Limit == 0 {
		paging.Limit = 100
	}
	return uc.Repo.FindActiveAccounts(account, paging)
}

func (uc *accountUseCase) Update(account *account.Account, hashPassword bool) error {
	if hashPassword {
		hashedPass, err := uc.HashPassword(account.Password)
		if err != nil {
			return err
		}
		account.Password = hashedPass
	}
	return uc.Repo.Update(account)
}
