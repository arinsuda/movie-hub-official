package user_stats_module

import "gorm.io/gorm"

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) GetUserStats(userID uint) (*UserStatsResponse, error) {
	row, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return toResponse(row), nil
}

func (s *Service) AddExperience(userID uint, expDelta int) error {
	return s.repo.UpsertStatus(userID, expDelta)
}

func (s *Service) GetLevel(userID uint) int {
	level, err := s.repo.GetLevelByUserID(userID)
	if err != nil {
		return 1
	}
	return level
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
