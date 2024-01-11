package bandrepo

import (
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"gorm.io/gorm"
)

type MemberRepo interface {
	Create(*band.Member) error
}

type memberRepo struct {
	db *gorm.DB
}

func NewMemberRepo(db *gorm.DB) MemberRepo {
	return &memberRepo{
		db: db,
	}
}

func (repo *memberRepo) Create(member *band.Member) error {
	return repo.db.Create(member).Error
}
