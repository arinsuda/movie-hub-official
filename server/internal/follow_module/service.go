package follow_module

import (
	"context"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	users "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
)

type Service struct {
	repo       *repository
	db         *gorm.DB
	achieveSvc achievementsmodule.Service
	notifSvc   *notification_module.Service
}

func NewService(
	db *gorm.DB,
	achieve achievementsmodule.Service,
	notif *notification_module.Service,
) *Service {
	return &Service{
		repo:       newRepository(db),
		db:         db,
		achieveSvc: achieve,
		notifSvc:   notif,
	}
}

func (s *Service) Follow(requesterID, targetID uint) (*FollowResponse, error) {
	follow, err := s.repo.Follow(requesterID, targetID)
	if err != nil {
		return nil, err
	}

	if follow.Status == StatusAccepted {
		s.onFollowAccepted(requesterID, targetID)
	}

	return &FollowResponse{
		FollowerID: follow.FollowerID,
		FolloweeID: follow.FolloweeID,
		Status:     string(follow.Status),
	}, nil
}

func (s *Service) Unfollow(requesterID, targetID uint) error {
	return s.repo.Unfollow(requesterID, targetID)
}

func (s *Service) AcceptFollow(requesterID, followerID uint) error {
	if err := s.repo.AcceptFollow(followerID, requesterID); err != nil {
		return err
	}

	s.onFollowAccepted(followerID, requesterID)

	return nil
}

func (s *Service) RejectFollow(requesterID, followerID uint) error {
	return s.repo.RejectFollow(followerID, requesterID)
}

func (s *Service) GetFollowers(userID uint) ([]UserSummary, error) {
	rows, err := s.repo.GetFollowers(userID)
	if err != nil {
		return nil, err
	}
	return toSummaryList(rows), nil
}

func (s *Service) GetFollowing(userID uint) ([]UserSummary, error) {
	rows, err := s.repo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}
	return toSummaryList(rows), nil
}

func (s *Service) GetPendingRequests(requesterID, userID uint) ([]UserSummary, error) {
	if requesterID != userID {
		return nil, ErrForbidden
	}
	rows, err := s.repo.GetPendingRequests(userID)
	if err != nil {
		return nil, err
	}
	return toSummaryList(rows), nil
}

func (s *Service) onFollowAccepted(followerID, followeeID uint) {
	ctx := context.Background()

	var followingCount int64
	s.db.Model(&UserFollow{}).
		Where("follower_id = ? AND status = ?", followerID, StatusAccepted).
		Count(&followingCount)
	_, _ = s.achieveSvc.Track(followerID, "following_count", int(followingCount))

	var followerCount int64
	s.db.Model(&UserFollow{}).
		Where("followee_id = ? AND status = ?", followeeID, StatusAccepted).
		Count(&followerCount)
	_, _ = s.achieveSvc.Track(followeeID, "follower_count", int(followerCount))

	if s.notifSvc != nil {
		if actor, err := s.getUserSummary(followerID); err == nil {
			_ = s.notifSvc.PushFollowedYou(ctx, followeeID, followerID, actor.Username)
		}
	}
}

func (s *Service) getUserSummary(userID uint) (*users.User, error) {
	var u users.User
	err := s.db.First(&u, userID).Error
	return &u, err
}

func toSummaryList(rows []listRow) []UserSummary {
	result := make([]UserSummary, len(rows))
	for i, r := range rows {
		result[i] = UserSummary{
			ID:          r.ID,
			Username:    r.Username,
			DisplayName: r.DisplayName,
			AvatarURL:   r.AvatarURL,
			Status:      r.Status,
		}
	}
	return result
}
