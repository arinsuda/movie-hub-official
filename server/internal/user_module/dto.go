package user_module

import "time"

type UserProfileResponse struct {
	ID             uint       `json:"id"`
	Username       string     `json:"username"`
	DisplayName    *string    `json:"display_name"`
	Bio            *string    `json:"bio"`
	AvatarURL      *string    `json:"avatar_url"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Gender         GenderType `json:"gender"`
	FavoriteGenres *string    `json:"favorite_genres"`
	ReviewCount    int        `json:"review_count"`
	FollowerCount  int        `json:"follower_count"`
	FollowingCount int        `json:"following_count"`
	IsPrivate      bool       `json:"is_private"`
	Role           string     `json:"role"`
}

type UpdateProfileRequest struct {
	DisplayName    *string    `json:"display_name"`
	Bio            *string    `json:"bio"`
	AvatarURL      *string    `json:"avatar_url"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Gender         GenderType `json:"gender"`
	GenderOther    *string    `json:"gender_other"`
	FavoriteGenres *string    `json:"favorite_genres"`
	IsPrivate      *bool      `json:"is_private"`
}
