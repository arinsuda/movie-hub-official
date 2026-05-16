package user_module

import (
	"gorm.io/gorm"
)

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) GetProfile(targetUserID uint, requesterID uint) (*UserProfileResponse, error) {
	user, reviewCount, followerCount, followingCount, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	profile := toProfileResponse(user, reviewCount, followerCount, followingCount)

	// ซ่อน fields ส่วนตัวถ้า profile เป็น private และไม่ใช่เจ้าของ
	if user.IsPrivate && requesterID != targetUserID {
		profile.Bio = nil
		profile.FavoriteGenres = nil
		profile.Gender = ""
		profile.DateOfBirth = nil
	}

	return profile, nil
}

func (s *Service) UpdateProfile(targetUserID uint, requesterID uint, req UpdateProfileRequest) (*UserProfileResponse, error) {
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
	if req.AvatarURL != nil {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.GenderOther != nil {
		updates["gender_other"] = req.GenderOther
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = req.DateOfBirth
	}
	if req.FavoriteGenres != nil {
		updates["favorite_genres"] = req.FavoriteGenres
	}
	if req.IsPrivate != nil {
		updates["is_private"] = req.IsPrivate
	}

	if len(updates) == 0 {
		return s.GetProfile(targetUserID, requesterID)
	}

	if err := s.repo.UpdateProfile(targetUserID, updates); err != nil {
		return nil, err
	}

	return s.GetProfile(targetUserID, requesterID)
}

func (s *Service) DeleteUser(targetUserID uint, requesterID uint, requesterRole string) error {
	if targetUserID != requesterID && requesterRole != "admin" {
		return ErrForbidden
	}
	return s.repo.DeleteUser(targetUserID)
}

// toProfileResponse รับ stats แยกจาก view แทนที่จะอ่านจาก User model โดยตรง
func toProfileResponse(u *User, reviewCount, followerCount, followingCount int) *UserProfileResponse {
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
		Role:           string(u.Role.RoleName),
	}
}
