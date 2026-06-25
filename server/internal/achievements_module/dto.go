package achievementsmodule

import "time"

type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p *PaginationQuery) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 20
	}
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func newPaginationMeta(page, limit int, total int64) PaginationMeta {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

type AchievementResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ExpReward   int    `json:"exp_reward"`
	ActionType  string `json:"action_type"`
	TargetCount int    `json:"target_count"`
}

type ListAchievementsResponse struct {
	Data       []AchievementResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

type UserAchievementResponse struct {
	AchievementID uint                `json:"achievement_id"`
	Achievement   AchievementResponse `json:"achievement"`
	CurrentCount  int                 `json:"current_count"`
	IsUnlocked    bool                `json:"is_unlocked"`
	UnlockedAt    *time.Time          `json:"unlocked_at,omitempty"`

	ProgressPct float64 `json:"progress_pct"`
}

type ListUserAchievementsResponse struct {
	Data       []UserAchievementResponse `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}

type UserAchievementFilter struct {
	PaginationQuery

	Unlocked   *bool  `query:"unlocked"`
	ActionType string `query:"action_type"`
}

type NewlyUnlocked struct {
	Achievement AchievementResponse `json:"achievement"`
	ExpGained   int                 `json:"exp_gained"`
}
