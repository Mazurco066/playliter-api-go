package concert

import "gorm.io/gorm"

type Concert struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
}
