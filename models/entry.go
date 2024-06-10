package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	Url   string
	Short string `gorm:"uniqueIndex:idx_short"`
	Hits  uint   `gorm:"index:idx_hits,sort:desc"`
	Tag   string
}

func (e Entry) String() string {
	return fmt.Sprintf("Made a shorty for %s - shorted with %s - with tag of: %s", e.Url, e.Short, e.Tag)
}
