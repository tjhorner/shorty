package main

import (
	"github.com/jinzhu/gorm"
)

type Link struct {
	gorm.Model `json:"-"`
	Suffix     string `json:"suffix"`
	URL        string `json:"url"`
	CreatorIP  string `json:"-"`
}
