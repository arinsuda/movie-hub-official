package user_module

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/arinsuda/movie-hub/internal/shared/storage"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"gorm.io/gorm"
)

type Service struct {
	repo     *repository
	minio    *storage.MinIOClient
	statsSvc *stats.Service // used only to read Level; no circular EXP dependency
}

func NewService(db *gorm.DB, mc *storage.MinIOClient, statsSvc *stats.Service) *Service {
	return &Service{repo: newRepository(db), minio: mc, statsSvc: statsSvc}
}

// ── GetProfile ────────────────────────────────────────────────────

func (s *Service) GetProfile(targetUserID, requesterID uint) (*UserProfileResponse, error) {
	user, reviewCount, followerCount, followingCount, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	level := s.fetchLevel(targetUserID)
	profile := toProfileResponse(user, reviewCount, followerCount, followingCount, level)

	if user.IsPrivate && requesterID != targetUserID {
		profile.Bio = nil
		profile.FavoriteGenres = nil
		profile.Gender = ""
		profile.DateOfBirth = nil
	}

	if profile.AvatarURL != nil && strings.HasPrefix(*profile.AvatarURL, "avatars/") {
		if url, err := s.minio.PresignURL(context.Background(), *profile.AvatarURL); err == nil {
			profile.AvatarURL = &url
		}
	}

	return profile, nil
}

// ── UpdateProfile ─────────────────────────────────────────────────

func (s *Service) UpdateProfile(
	targetUserID, requesterID uint,
	req UpdateProfileRequest,
	avatarFile multipart.File,
	avatarHeader *multipart.FileHeader,
) (*UserProfileResponse, error) {
	if targetUserID != requesterID {
		return nil, ErrForbidden
	}

	updates := map[string]any{}
	if req.DisplayName != nil {
		updates["display_name"] = req.DisplayName
	}
	if req.Bio != nil {
		updates["bio"] = req.Bio
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = req.DateOfBirth
	}
	if req.Gender != "" {
		if !isValidGender(req.Gender) {
			return nil, errors.New("invalid gender")
		}
		updates["gender"] = req.Gender
	}
	if req.GenderOther != nil {
		updates["gender_other"] = req.GenderOther
	}
	if req.FavoriteGenres != nil {
		updates["favorite_genres"] = req.FavoriteGenres
	}
	if req.IsPrivate != nil {
		updates["is_private"] = req.IsPrivate
	}

	if avatarFile != nil && avatarHeader != nil {
		oldKey := s.currentAvatarKey(targetUserID)
		newKey, err := s.minio.UploadAvatar(context.Background(), targetUserID, avatarFile, avatarHeader)
		if err != nil {
			return nil, err
		}
		updates["avatar_url"] = newKey
		if oldKey != "" {
			_ = s.minio.DeleteObject(context.Background(), oldKey)
		}
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateProfile(targetUserID, updates); err != nil {
			return nil, err
		}
	}

	return s.GetProfile(targetUserID, requesterID)
}

// ── DeleteUser ────────────────────────────────────────────────────

func (s *Service) DeleteUser(targetUserID, requesterID uint, requesterRole string) error {
	if targetUserID != requesterID && requesterRole != "admin" {
		return ErrForbidden
	}
	return s.repo.DeleteUser(targetUserID)
}

// ── UpdateFavoriteGenres ──────────────────────────────────────────

func (s *Service) UpdateFavoriteGenres(targetUserID, requesterID uint, genres []int) (*UserProfileResponse, error) {
	if targetUserID != requesterID {
		return nil, ErrForbidden
	}

	var stored *string
	if len(genres) > 0 {
		b, err := json.Marshal(genres)
		if err != nil {
			return nil, err
		}
		str := string(b)
		stored = &str
	}

	if err := s.repo.UpdateProfile(targetUserID, map[string]any{"favorite_genres": stored}); err != nil {
		return nil, err
	}
	return s.GetProfile(targetUserID, requesterID)
}

// ── Helpers ───────────────────────────────────────────────────────

// fetchLevel reads the user's current level from user_stats_module.
// Returns 1 on any error so the profile response is never broken.
func (s *Service) fetchLevel(userID uint) int {
	if s.statsSvc == nil {
		return 1
	}
	st, err := s.statsSvc.GetUserStats(userID)
	if err != nil || st == nil {
		return 1
	}
	return st.Level
}

func (s *Service) currentAvatarKey(userID uint) string {
	user, _, _, _, err := s.repo.FindByID(userID)
	if err != nil || user.AvatarURL == nil {
		return ""
	}
	if strings.HasPrefix(*user.AvatarURL, "avatars/") {
		return *user.AvatarURL
	}
	return ""
}

func toProfileResponse(u *User, reviewCount, followerCount, followingCount, level int) *UserProfileResponse {
	return &UserProfileResponse{
		ID:             u.ID,
		Username:       u.Username,
		DisplayName:    u.DisplayName,
		Bio:            u.Bio,
		AvatarURL:      u.AvatarURL,
		Gender:         u.Gender,
		FavoriteGenres: u.FavoriteGenres,
		DateOfBirth:    u.DateOfBirth,
		ReviewCount:    reviewCount,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		IsPrivate:      u.IsPrivate,
		Level:          level,
		Role:           string(u.Role.RoleName),
		CreatedAt:      u.CreatedAt,
	}
}
