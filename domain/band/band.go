package band

import "gorm.io/gorm"

type Band struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
}
