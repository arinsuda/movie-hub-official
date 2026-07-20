package bmol_module

import (
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(item *BMOLItem) error {
	return r.db.Create(item).Error
}

func (r *repository) Update(id uint, rank int) error {
	return r.db.Model(&BMOLItem{}).Where("id = ?", id).Update("rank", rank).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&BMOLItem{}, id).Error
}

func (r *repository) FindOne(id uint) (*BMOLItem, error) {
	var item BMOLItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindByUser(userID uint, mediaType *movie_module.MediaType) ([]BMOLItem, error) {
	var items []BMOLItem
	query := r.db.Where("user_id = ? AND deleted_at IS NULL", userID)
	if mediaType != nil && *mediaType != "" {
		query = query.Where("media_type = ?", *mediaType)
	}
	err := query.Order("rank asc").Order("created_at asc").Find(&items).Error
	return items, err
}

func (r *repository) FindDuplicate(userID uint, mediaID int, mediaType movie_module.MediaType) (*BMOLItem, error) {
	var item BMOLItem
	err := r.db.Where("user_id = ? AND media_id = ? AND media_type = ? AND deleted_at IS NULL", userID, mediaID, mediaType).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
