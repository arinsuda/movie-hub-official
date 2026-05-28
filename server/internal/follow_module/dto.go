package follow_module

// FollowResponse ส่งกลับหลัง follow สำเร็จ
// status บอกให้ client รู้ว่า follow ทันที หรือรอ approve
type FollowResponse struct {
	FollowerID uint   `json:"follower_id"`
	FolloweeID uint   `json:"followee_id"`
	Status     string `json:"status"` // "accepted" | "pending"
}

// UserSummary ใช้ใน followers / following / pending lists
type UserSummary struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
	Status      string  `json:"status"`
}
