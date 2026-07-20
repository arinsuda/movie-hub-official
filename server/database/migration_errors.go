package database

type DuplicateReviewGroup struct {
	UserID    uint   `json:"user_id"`
	MediaID   int    `json:"media_id"`
	MediaType string `json:"media_type"`
	ReviewIDs []uint `json:"review_ids"`
}

type ErrActiveDuplicateReviews struct {
	Groups []DuplicateReviewGroup `json:"groups"`
}

func (e *ErrActiveDuplicateReviews) Error() string {
	return "active duplicate reviews detected"
}
