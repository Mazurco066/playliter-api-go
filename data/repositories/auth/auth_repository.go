package authrepo

import "gorm.io/gorm"

type Repo interface {
}

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) Repo {
	return &AuthRepo{
		db: db,
	}
}
