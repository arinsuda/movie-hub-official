package library_module

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) AddItem(userID uint, req AddItemRequest) (*LibraryItemResponse, error) {
	if req.MediaType != MediaMovie && req.MediaType != MediaSeries {
		return nil, ErrInvalidMediaType 
	}
	if req.ListType != ListWatchlist && req.ListType != ListFavorite && req.ListType != ListWatched {
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

	if len(req.Tags) > 0 {
		b, _ := json.Marshal(req.Tags)
		s := string(b)
		item.Tags = &s
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return toResponse(item), nil
}

func (s *Service) GetLibrary(userID uint, listType *ListType, mediaType *MediaType) ([]LibraryItemResponse, error) {
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

func (s *Service) GetMediaStatus(userID uint, mediaID int, mediaType MediaType) (*MediaStatusResponse, error) {
	items, err := s.repo.FindMediaStatus(userID, mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	inLists := make([]ListType, len(items))
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
	if item.Tags != nil {
		_ = json.Unmarshal([]byte(*item.Tags), &tags)
	}
	if tags == nil {
		tags = []string{}
	}

	return &LibraryItemResponse{
		ID:        item.ID,
		MediaID:   item.MediaID,
		MediaType: item.MediaType,
		ListType:  item.ListType,
		WatchedAt: item.WatchedAt,
		Tags:      tags,
		Note:      item.Note,
		CreatedAt: item.CreatedAt,
	}
}
