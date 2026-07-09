package follow_module

import (
	"context"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/shared"
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

	switch follow.Status {
	case StatusAccepted:
		s.onFollowEstablished(requesterID, targetID)
		s.notifyNewFollower(requesterID, targetID)
	case StatusPending:
		s.notifyFollowRequested(requesterID, targetID)
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

	s.onFollowEstablished(followerID, requesterID)
	s.notifyFollowAccepted(requesterID, followerID)

	return nil
}

func (s *Service) RejectFollow(requesterID, followerID uint) error {
	return s.repo.RejectFollow(followerID, requesterID)
}

func (s *Service) GetFollowers(requesterID, userID uint) ([]UserSummary, error) {
	canView, err := s.repo.canViewFollowList(requesterID, userID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, ErrForbidden
	}

	rows, err := s.repo.GetFollowers(userID)
	if err != nil {
		return nil, err
	}
	return toSummaryList(rows), nil
}

func (s *Service) GetFollowing(requesterID, userID uint) ([]UserSummary, error) {
	canView, err := s.repo.canViewFollowList(requesterID, userID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, ErrForbidden
	}

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

func (s *Service) GetRelationshipStatus(requesterID, targetID uint) (*RelationshipStatus, error) {
	outbound, err := s.repo.GetStatus(requesterID, targetID)
	if err != nil {
		return nil, err
	}
	inbound, err := s.repo.GetStatus(targetID, requesterID)
	if err != nil {
		return nil, err
	}

	status := &RelationshipStatus{}
	if outbound != nil {
		status.IsFollowing = outbound.Status == StatusAccepted
		status.FollowStatus = string(outbound.Status)
	}
	if inbound != nil {
		status.IsFollowedBy = inbound.Status == StatusAccepted
	}
	return status, nil
}

func (s *Service) onFollowEstablished(followerID, followeeID uint) {
	ctx := context.Background()

	var followingCount int64
	s.db.Model(&UserFollow{}).
		Where("follower_id = ? AND status = ?", followerID, StatusAccepted).
		Count(&followingCount)
	shared.TrackAndNotify(ctx, s.achieveSvc, s.notifSvc, followerID, "following_count", int(followingCount))

	var followerCount int64
	s.db.Model(&UserFollow{}).
		Where("followee_id = ? AND status = ?", followeeID, StatusAccepted).
		Count(&followerCount)
	shared.TrackAndNotify(ctx, s.achieveSvc, s.notifSvc, followeeID, "follower_count", int(followerCount))
}

func (s *Service) notifyNewFollower(followerID, followeeID uint) {
	if s.notifSvc == nil {
		return
	}
	actor, err := s.getUserSummary(followerID)
	if err != nil {
		return
	}
	_ = s.notifSvc.PushFollowedYou(context.Background(), followeeID, followerID, actor.Username)
}

func (s *Service) notifyFollowRequested(followerID, followeeID uint) {
	if s.notifSvc == nil {
		return
	}
	actor, err := s.getUserSummary(followerID)
	if err != nil {
		return
	}
	_ = s.notifSvc.PushFollowRequested(context.Background(), followeeID, followerID, actor.Username)
}

func (s *Service) notifyFollowAccepted(accepterID, followerID uint) {
	if s.notifSvc == nil {
		return
	}
	actor, err := s.getUserSummary(accepterID)
	if err != nil {
		return
	}
	_ = s.notifSvc.PushFollowAccepted(context.Background(), followerID, accepterID, actor.Username)
}

func (s *Service) getUserSummary(userID uint) (*users.User, error) {
	var u users.User
	err := s.db.First(&u, userID).Error
	return &u, err
}

func (s *Service) GetFollowStats(requesterID, targetID uint) (*FollowStatsResponse, error) {
	followers, err := s.repo.CountFollowers(targetID)
	if err != nil {
		return nil, err
	}
	following, err := s.repo.CountFollowing(targetID)
	if err != nil {
		return nil, err
	}

	isFollowing := false
	if requesterID != targetID {
		status, err := s.repo.GetStatus(requesterID, targetID)
		if err != nil {
			return nil, err
		}
		isFollowing = status != nil && status.Status == StatusAccepted
	}

	return &FollowStatsResponse{
		UserID:      targetID,
		Followers:   followers,
		Following:   following,
		IsFollowing: isFollowing,
	}, nil
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
