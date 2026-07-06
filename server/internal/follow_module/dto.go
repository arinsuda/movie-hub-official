package follow_module

type FollowResponse struct {
	FollowerID uint   `json:"follower_id"`
	FolloweeID uint   `json:"followee_id"`
	Status     string `json:"status"`
}

type UserSummary struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
	Status      string  `json:"status"`
}
