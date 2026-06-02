package statsmodule

import (
	"errors"

	"github.com/arinsuda/movie-hub/internal/library_module"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrStatNotFound = errors.New("stat not found")

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetStats(mediaID int, mediaType library_module.MediaType) (*MediaStat, error) {
	var stat MediaStat
	err := r.db.Where("media_id = ? AND media_type = ?", mediaID, mediaType).First(&stat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStatNotFound
	}
	return &stat, err
}

func (r *repository) Upsert(mediaID int, mediaType library_module.MediaType) error {
	stat := MediaStat{
		MediaID:   mediaID,
		MediaType: mediaType,
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&stat).Error
}

func (r *repository) Increment(mediaID int, mediaType library_module.MediaType, field IncrementField, delta int) error {
	result := r.db.Model(&MediaStat{}).
		Where("media_id = ? AND media_type = ?", mediaID, mediaType).
		UpdateColumn(string(field), gorm.Expr(string(field)+" + ?", delta))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStatNotFound
	}
	return nil
}
