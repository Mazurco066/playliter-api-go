package accountrepo

import (
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*account.Account) error
	FindActiveAccounts(*account.Account, *commoninputs.PagingParams) ([]*account.Account, error)
	FindByUserEmail(email string) (*account.Account, error)
	FindByUsernameOrEmail(filter string) (*account.Account, error)
	Update(*account.Account) error
}

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) Repo {
	return &AccountRepo{
		db: db,
	}
}

func (repo *AccountRepo) Create(account *account.Account) error {
	return repo.db.Create(account).Error
}

func (repo *AccountRepo) FindActiveAccounts(a *account.Account, p *commoninputs.PagingParams) ([]*account.Account, error) {
	var results []*account.Account
	if err := repo.db.Where(
		"is_active = ? AND id != ?",
		true, a.ID,
	).
		Limit(p.Limit).
		Offset(p.Offset).
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *AccountRepo) FindByUserEmail(email string) (*account.Account, error) {
	var account account.Account
	if err := repo.db.Where("email = ? AND is_active = ?", email, true).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *AccountRepo) FindByUsernameOrEmail(filter string) (*account.Account, error) {
	var account account.Account
	if err := repo.db.Where(
		"email = ? OR username = ? AND is_active = ?",
		filter, filter, true,
	).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *AccountRepo) Update(acount *account.Account) error {
	return repo.db.Save(acount).Error
}
