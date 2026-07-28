package watch_log_module

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrWatchLogNotFound  = errors.New("watch log not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidDate       = errors.New("invalid date format")
	ErrFutureDate        = errors.New("cannot set watch date in the future")
	ErrInvalidVisibility = errors.New("invalid visibility option")
	ErrInvalidMediaType  = errors.New("invalid media type")
	ErrInvalidMediaID    = errors.New("invalid media id")
)

type Repository interface {
	Create(log *MediaWatchLog) error
	FindByID(id uint) (*MediaWatchLog, error)
	FindByUserAndMedia(userID uint, mediaType string, mediaID int) ([]MediaWatchLog, error)
	FindByUser(userID uint, pq PaginationQuery) ([]MediaWatchLog, int64, error)
	Update(id uint, updates map[string]any) error
	Delete(id uint) error
	ComputeSummary(userID uint, mediaType string, mediaID int) (WatchSummary, error)
	CountByUserAndMedia(userID uint, mediaType string, mediaID int) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(log *MediaWatchLog) error {
	return r.db.Create(log).Error
}

func (r *repository) FindByID(id uint) (*MediaWatchLog, error) {
	var log MediaWatchLog
	err := r.db.First(&log, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *repository) FindByUserAndMedia(userID uint, mediaType string, mediaID int) ([]MediaWatchLog, error) {
	var logs []MediaWatchLog
	err := r.db.Where("user_id = ? AND media_type = ? AND media_id = ?", userID, mediaType, mediaID).
		Order("watched_on DESC, id DESC").Find(&logs).Error
	return logs, err
}

func (r *repository) FindByUser(userID uint, pq PaginationQuery) ([]MediaWatchLog, int64, error) {
	var logs []MediaWatchLog
	var total int64

	query := r.db.Model(&MediaWatchLog{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("watched_on DESC, id DESC").
		Offset(pq.Offset()).
		Limit(pq.Limit).
		Find(&logs).Error

	return logs, total, err
}

func (r *repository) Update(id uint, updates map[string]any) error {
	return r.db.Model(&MediaWatchLog{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&MediaWatchLog{}, id).Error
}

func (r *repository) ComputeSummary(userID uint, mediaType string, mediaID int) (WatchSummary, error) {
	var count int64
	var minDate *string
	var maxDate *string

	err := r.db.Model(&MediaWatchLog{}).
		Where("user_id = ? AND media_type = ? AND media_id = ?", userID, mediaType, mediaID).
		Select("COUNT(*), MIN(CAST(watched_on AS TEXT)), MAX(CAST(watched_on AS TEXT))").
		Row().
		Scan(&count, &minDate, &maxDate)
	if err != nil {
		return WatchSummary{}, err
	}

	return WatchSummary{
		WatchCount:     int(count),
		FirstWatchedOn: minDate,
		LastWatchedOn:  maxDate,
		HasWatched:     count > 0,
	}, nil
}

func (r *repository) CountByUserAndMedia(userID uint, mediaType string, mediaID int) (int64, error) {
	var count int64
	err := r.db.Model(&MediaWatchLog{}).
		Where("user_id = ? AND media_type = ? AND media_id = ?", userID, mediaType, mediaID).
		Count(&count).Error
	return count, err
}
