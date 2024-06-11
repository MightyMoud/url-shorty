package services

import (
	"context"
	"errors"

	"github.com/ms-mousa/url-shorty/models"
	"gorm.io/gorm"
)

func FindEntry(ctx context.Context, queryEntry *models.Entry) (models.Entry, error) {
	db := ctx.Value("db").(*gorm.DB)
	if res := db.Find(&queryEntry, &queryEntry); res.Error != nil {
		return *queryEntry, errors.New("Unfound")
	}
	queryEntry.Hits = queryEntry.Hits + 1
	db.Save(&queryEntry)
	return *queryEntry, nil
}

func AddEntry(ctx context.Context, destEntry *models.Entry) (models.Entry, error) {
	db := ctx.Value("db").(*gorm.DB)

	if result := db.Create(&destEntry); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			if res := db.First(&destEntry); res.Error != nil {
				destEntry.Hits = destEntry.Hits + 1
				db.Save(&destEntry)
				return *destEntry, nil
			}
			return *destEntry, nil
		}
		return *destEntry, errors.New("Something Went Wrong bro~")
	}
	return *destEntry, nil
}
