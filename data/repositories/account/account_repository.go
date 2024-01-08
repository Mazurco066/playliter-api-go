package accountrepo

import (
	"github.com/mazurco066/playliter-api-go/domain/account"
	"gorm.io/gorm"
)

type Repo interface {
	FindByUsernameOrEmail(filter string) (*account.Account, error)
}

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) Repo {
	return &AccountRepo{
		db: db,
	}
}

func (repo *AccountRepo) FindByUsernameOrEmail(filter string) (*account.Account, error) {
	var account account.Account
	if err := repo.db.Where(
		"email = ? OR username = ?",
		filter, filter,
	).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
