package models

import "gorm.io/gorm"

type StackDraft struct {
	gorm.Model
	Name string `gorm:"unique"`
	Data string `gorm:"type:text"`
}
