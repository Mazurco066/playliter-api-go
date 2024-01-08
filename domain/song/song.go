package song

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
}
