package band

import "gorm.io/gorm"

type BandRequest struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
}
