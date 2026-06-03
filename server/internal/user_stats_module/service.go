package user_stats_module

import "gorm.io/gorm"

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) GetUserStats(userID uint) (*UserStatsResponse, error) {
	stat, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return toResponse(stat), nil
}

func toResponse(s *UserStat) *UserStatsResponse {
	return &UserStatsResponse{
		UserID:         s.UserID,
		ReviewCount:    s.ReviewCount,
		LikeCount:      s.LikeCount,
		WatchlistCount: s.WatchlistCount,
		WatchedCount:   s.WatchedCount,
		FollowerCount:  s.FollowerCount,
		FollowingCount: s.FollowingCount,
	}
}
