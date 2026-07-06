package notification_module

import (
	"context"
	"fmt"

	user_module "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
)

type UserProvider interface {
	FindByID(id uint) (*user_module.User, int, int, int, error)
	FindFollowerIDs(userID uint) ([]uint, error)
}

type Service struct {
	repo         *repository
	userProvider UserProvider
	hub          *Hub
}

func NewService(db *gorm.DB, userProvider UserProvider, hub *Hub) *Service {
	return &Service{
		repo:         newRepository(db),
		userProvider: userProvider,
		hub:          hub,
	}
}

func (s *Service) ListNotifications(ctx context.Context, userID uint, q ListNotificationsQuery) (*NotificationListResponse, error) {
	rows, total, err := s.repo.FindByUser(ctx, userID, q)
	if err != nil {
		return nil, err
	}

	unread, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return nil, err
	}

	page, pageSize := normalizePagination(q.Page, q.PageSize)

	responses := make([]NotificationResponse, 0, len(rows))
	for _, n := range rows {
		responses = append(responses, s.toResponse(n))
	}

	return &NotificationListResponse{
		Notifications: responses,
		UnreadCount:   unread,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

func (s *Service) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.repo.CountUnread(ctx, userID)
}

func (s *Service) MarkRead(ctx context.Context, userID uint, req MarkReadRequest) error {
	if err := s.repo.MarkRead(ctx, userID, req.IDs); err != nil {
		return err
	}
	s.hub.EmitRead(userID, req.IDs)
	if count, err := s.repo.CountUnread(ctx, userID); err == nil {
		s.hub.EmitUnreadCount(userID, count)
	}
	return nil
}

func (s *Service) DeleteNotifications(ctx context.Context, userID uint, ids []uint) error {
	if err := s.repo.DeleteByUser(ctx, userID, ids); err != nil {
		return err
	}
	s.hub.EmitDeleted(userID, ids)
	return nil
}

func (s *Service) createAndEmit(ctx context.Context, ns []Notification) error {
	if len(ns) == 0 {
		return nil
	}
	if err := s.repo.CreateBatch(ctx, ns); err != nil {
		return err
	}
	for _, n := range ns {
		s.hub.EmitNew(n.UserID, s.toResponse(n))
	}
	return nil
}

func (s *Service) PushFollowedYou(ctx context.Context, targetUserID, actorID uint, actorUsername string) error {
	n := Notification{
		UserID:    targetUserID,
		ActorID:   &actorID,
		Type:      NotifFollowedYou,
		TargetID:  &actorID,
		TargetRef: ptr("user"),
		Message:   fmt.Sprintf("%s started following you", actorUsername),
	}
	return s.createAndEmit(ctx, []Notification{n})
}

func (s *Service) PushMovieNowPlaying(ctx context.Context, userIDs []uint, movieID uint, movieTitle string) error {
	ns := make([]Notification, 0, len(userIDs))
	for _, uid := range userIDs {
		ns = append(ns, Notification{
			UserID:    uid,
			ActorID:   nil,
			Type:      NotifMovieNowPlaying,
			TargetID:  &movieID,
			TargetRef: ptr("movie"),
			Message:   fmt.Sprintf(`"%s" is now playing — time to watch!`, movieTitle),
		})
	}
	return s.createAndEmit(ctx, ns)
}

func (s *Service) PushFollowingActivity(
	ctx context.Context,
	actorID uint,
	notifType NotificationType,
	targetID *uint,
	targetRef *string,
	message string,
) error {
	followerIDs, err := s.userProvider.FindFollowerIDs(actorID)
	if err != nil {
		return err
	}

	ns := make([]Notification, 0, len(followerIDs))
	for _, fid := range followerIDs {
		ns = append(ns, Notification{
			UserID:    fid,
			ActorID:   &actorID,
			Type:      notifType,
			TargetID:  targetID,
			TargetRef: targetRef,
			Message:   message,
		})
	}
	return s.createAndEmit(ctx, ns)
}

func (s *Service) PushFollowingReviewed(ctx context.Context, actorID uint, actorUsername string, reviewID uint, movieTitle string) error {
	ref := "review"
	msg := fmt.Sprintf("%s reviewed %q", actorUsername, movieTitle)
	return s.PushFollowingActivity(ctx, actorID, NotifFollowingReviewed, &reviewID, &ref, msg)
}

func (s *Service) PushFollowingLikedReview(ctx context.Context, actorID uint, actorUsername string, reviewID uint, movieTitle string) error {
	ref := "review"
	msg := fmt.Sprintf("%s liked a review of %q", actorUsername, movieTitle)
	return s.PushFollowingActivity(ctx, actorID, NotifFollowingLikedReview, &reviewID, &ref, msg)
}

func (s *Service) PushFollowingAddedWatchlist(ctx context.Context, actorID uint, actorUsername string, movieID uint, movieTitle string) error {
	ref := "movie"
	msg := fmt.Sprintf("%s added %q to their watchlist", actorUsername, movieTitle)
	return s.PushFollowingActivity(ctx, actorID, NotifFollowingAddedWatchlist, &movieID, &ref, msg)
}

func (s *Service) PushFollowingAddedWatched(ctx context.Context, actorID uint, actorUsername string, movieID uint, movieTitle string) error {
	ref := "movie"
	msg := fmt.Sprintf("%s marked %q as watched", actorUsername, movieTitle)
	return s.PushFollowingActivity(ctx, actorID, NotifFollowingAddedWatched, &movieID, &ref, msg)
}

func (s *Service) toResponse(n Notification) NotificationResponse {
	resp := NotificationResponse{
		ID:        n.ID,
		Type:      n.Type,
		Message:   n.Message,
		IsRead:    n.IsRead,
		ReadAt:    n.ReadAt,
		TargetID:  n.TargetID,
		TargetRef: n.TargetRef,
		CreatedAt: n.CreatedAt,
	}

	if n.ActorID != nil {
		if actor, _, _, _, err := s.userProvider.FindByID(*n.ActorID); err == nil {
			resp.Actor = &ActorSummary{
				ID:          actor.ID,
				Username:    actor.Username,
				DisplayName: actor.DisplayName,
				AvatarURL:   actor.AvatarURL,
			}
		}
	}

	return resp
}

func ptr[T any](v T) *T { return &v }
