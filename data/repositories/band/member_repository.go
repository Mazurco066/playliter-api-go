package bandrepo

import (
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"gorm.io/gorm"
)

type MemberRepo interface {
	Create(*band.Member) error
	FindById(uint) (*band.Member, error)
	Remove(*band.Member) error
	Update(*band.Member) error
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

func (repo *memberRepo) FindById(id uint) (*band.Member, error) {
	var member band.Member
	if err := repo.db.
		Where("id = ?", id).
		Preload("Band").
		Preload("Account").
		First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (repo *memberRepo) Update(member *band.Member) error {
	return repo.db.Save(member).Error
}

func (repo *memberRepo) Remove(member *band.Member) error {
	return repo.db.Delete(member).Error
}
