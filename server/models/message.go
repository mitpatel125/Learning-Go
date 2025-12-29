package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Content string
	Read bool
}
