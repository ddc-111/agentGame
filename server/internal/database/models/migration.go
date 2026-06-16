package models

import "gorm.io/gorm"

type Migration struct {
	Version string
	Name    string
	Up      func(tx *gorm.DB) error
	Down    func(tx *gorm.DB) error
}
