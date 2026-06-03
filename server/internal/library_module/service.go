package library_module

import (
	"encoding/json"
	"time"

	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) AddItem(userID uint, req AddItemRequest) (*LibraryItemResponse, error) {
	if req.MediaType != movie_module.MediaMovie && req.MediaType != movie_module.MediaSeries {
		return nil, ErrInvalidMediaType 
	}
	if req.ListType != movie_module.ListWatchlist && req.ListType != movie_module.ListFavorite && req.ListType != movie_module.ListWatched {
		return nil, ErrInvalidListType
	}
	if req.MediaID <= 0 {
		return nil, ErrInvalidMediaID
	}
	item := &LibraryItem{
		UserID:    userID,
		MediaID:   req.MediaID,
		MediaType: req.MediaType,
		ListType:  req.ListType,
		Note:      req.Note,
	}

	if req.WatchedAt != nil {
		t, err := time.Parse("2006-01-02", *req.WatchedAt)
		if err != nil {
			return nil, ErrInvalidWatchedAt
		}
		item.WatchedAt = &t
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return toResponse(item), nil
}

func (s *Service) GetLibrary(userID uint, listType *movie_module.ListType, mediaType *movie_module.MediaType) ([]LibraryItemResponse, error) {
	items, err := s.repo.FindByUser(userID, listType, mediaType)
	if err != nil {
		return nil, err
	}

	responses := make([]LibraryItemResponse, len(items))
	for i, item := range items {
		responses[i] = *toResponse(&item)
	}
	return responses, nil
}

func (s *Service) RemoveItem(itemID, requesterID uint) error {
	item, err := s.repo.FindOne(itemID, requesterID)
	if err != nil {
		return err
	}
	if item.UserID != requesterID {
		return ErrForbidden
	}
	return s.repo.Delete(itemID)
}

func (s *Service) UpdateItem(itemID, requesterID uint, req UpdateItemRequest) (*LibraryItemResponse, error) {
	item, err := s.repo.FindOne(itemID, requesterID)
	if err != nil {
		return nil, err
	}
	if item.UserID != requesterID {
		return nil, ErrForbidden
	}

	updates := map[string]any{}

	if req.WatchedAt != nil {
		t, err := time.Parse("2006-01-02", *req.WatchedAt)
		if err != nil {
			return nil, ErrInvalidWatchedAt
		}
		updates["watched_at"] = t
	}

	if req.Tags != nil {
		b, _ := json.Marshal(req.Tags)
		updates["tags"] = string(b)
	}

	if req.Note != nil {
		updates["note"] = req.Note
	}

	if len(updates) == 0 {
		return toResponse(item), nil
	}

	if err := s.repo.Update(itemID, updates); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindOne(itemID, requesterID)
	if err != nil {
		return nil, err
	}
	return toResponse(updated), nil
}

func (s *Service) GetMediaStatus(userID uint, mediaID int, mediaType movie_module.MediaType) (*MediaStatusResponse, error) {
	items, err := s.repo.FindMediaStatus(userID, mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	inLists := make([]movie_module.ListType, len(items))
	for i, item := range items {
		inLists[i] = item.ListType
	}

	return &MediaStatusResponse{
		MediaID:   mediaID,
		MediaType: mediaType,
		InLists:   inLists,
	}, nil
}

func toResponse(item *LibraryItem) *LibraryItemResponse {
	var tags []string
	if tags == nil {
		tags = []string{}
	}

	return &LibraryItemResponse{
		ID:        item.ID,
		MediaID:   item.MediaID,
		MediaType: item.MediaType,
		ListType:  item.ListType,
		Tags:      tags,
		CreatedAt: item.CreatedAt,
	}
}
