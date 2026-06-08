package user_stats_module

import "gorm.io/gorm"

// Service implements ExpAdder. Other modules depend on the interface, not this struct.
type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

// GetUserStats returns the combined stats + level/exp for a user.
func (s *Service) GetUserStats(userID uint) (*UserStatsResponse, error) {
	row, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return toResponse(row), nil
}

// AddExperience satisfies the ExpAdder port.
// expDelta may be negative (e.g. review deleted).
func (s *Service) AddExperience(userID uint, expDelta int) error {
	return s.repo.UpsertStatus(userID, expDelta)
}

func toResponse(r *statRow) *UserStatsResponse {
	return &UserStatsResponse{
		UserID:         r.UserID,
		Level:          r.Level,
		CurrentExp:     r.CurrentExp,
		ReviewCount:    r.ReviewCount,
		LikeCount:      r.LikeCount,
		WatchlistCount: r.WatchlistCount,
		WatchedCount:   r.WatchedCount,
		FollowerCount:  r.FollowerCount,
		FollowingCount: r.FollowingCount,
	}
}
