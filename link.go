package main

import (
	"github.com/jinzhu/gorm"
)

type Link struct {
	gorm.Model `json:"-"`
	Suffix     string `json:"suffix" gorm:"UNIQUE"`
	URL        string `json:"url"`
	CreatorIP  string `json:"-"`
}
