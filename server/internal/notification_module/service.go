package notification_module

import (
	"context"
	"fmt"

	user_module "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
)

// UserProvider interface สำหรับดึงข้อมูล user (ตัดการ import circular)
type UserProvider interface {
	FindByID(id uint) (*user_module.User, int, int, int, error)
	FindFollowerIDs(userID uint) ([]uint, error) // ดึง id ของ followers ทั้งหมด
}

type Service struct {
	repo         *repository
	userProvider UserProvider
}

func NewService(db *gorm.DB, userProvider UserProvider) *Service {
	return &Service{
		repo:         newRepository(db),
		userProvider: userProvider,
	}
}

// ─── Read Operations ──────────────────────────────────────────────────────────

func (s *Service) ListNotifications(ctx context.Context, userID uint, q ListNotificationsQuery) (*NotificationListResponse, error) {
	rows, total, err := s.repo.FindByUser(userID, q)
	if err != nil {
		return nil, err
	}

	unread, err := s.repo.CountUnread(userID)
	if err != nil {
		return nil, err
	}

	page, pageSize := normalizePagination(q.Page, q.PageSize)

	responses := make([]NotificationResponse, 0, len(rows))
	for _, n := range rows {
		resp, err := s.toResponse(n)
		if err != nil {
			continue // skip ถ้า actor โดนลบ
		}
		responses = append(responses, resp)
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
	return s.repo.CountUnread(userID)
}

// ─── Mark / Delete ────────────────────────────────────────────────────────────

func (s *Service) MarkRead(ctx context.Context, userID uint, req MarkReadRequest) error {
	return s.repo.MarkRead(userID, req.IDs)
}

func (s *Service) DeleteNotifications(ctx context.Context, userID uint, ids []uint) error {
	return s.repo.DeleteByUser(userID, ids)
}

// ─── Push helpers (เรียกจาก modules อื่น) ────────────────────────────────────

// PushFollowedYou แจ้ง targetUser ว่ามีคน follow
func (s *Service) PushFollowedYou(ctx context.Context, targetUserID, actorID uint, actorUsername string) error {
	return s.repo.Create(&Notification{
		UserID:    targetUserID,
		ActorID:   &actorID,
		Type:      NotifFollowedYou,
		TargetID:  &actorID,
		TargetRef: ptr("user"),
		Message:   fmt.Sprintf("%s started following you", actorUsername),
	})
}

// PushMovieNowPlaying แจ้ง users ที่เพิ่มหนังใน watchlist ว่าหนังเข้าฉายแล้ว
// userIDs = รายชื่อ users ที่ต้องการแจ้ง, movieID & movieTitle = ข้อมูลหนัง
func (s *Service) PushMovieNowPlaying(ctx context.Context, userIDs []uint, movieID uint, movieTitle string) error {
	if len(userIDs) == 0 {
		return nil
	}

	ns := make([]Notification, 0, len(userIDs))
	for _, uid := range userIDs {
		ns = append(ns, Notification{
			UserID:    uid,
			ActorID:   nil, // system notification
			Type:      NotifMovieNowPlaying,
			TargetID:  &movieID,
			TargetRef: ptr("movie"),
			Message:   fmt.Sprintf(`"%s" is now playing — time to watch!`, movieTitle),
		})
	}
	return s.repo.CreateBatch(ns)
}

// PushFollowingActivity fan-out notification ไปหา followers ของ actorID
// เรียกได้จากหลาย event: review, like, watchlist
func (s *Service) PushFollowingActivity(
	ctx context.Context,
	actorID uint,
	actorUsername string,
	notifType NotificationType,
	targetID *uint,
	targetRef *string,
	message string,
) error {
	followerIDs, err := s.userProvider.FindFollowerIDs(actorID)
	if err != nil {
		return err
	}
	if len(followerIDs) == 0 {
		return nil
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
	return s.repo.CreateBatch(ns)
}

// PushFollowingReviewed shortcut สำหรับ review event
func (s *Service) PushFollowingReviewed(ctx context.Context, actorID uint, actorUsername string, reviewID uint, movieTitle string) error {
	rid := reviewID
	ref := "review"
	msg := fmt.Sprintf("%s reviewed %q", actorUsername, movieTitle)
	return s.PushFollowingActivity(ctx, actorID, actorUsername, NotifFollowingReviewed, &rid, &ref, msg)
}

// PushFollowingLikedReview shortcut สำหรับ like event
func (s *Service) PushFollowingLikedReview(ctx context.Context, actorID uint, actorUsername string, reviewID uint, movieTitle string) error {
	rid := reviewID
	ref := "review"
	msg := fmt.Sprintf("%s liked a review of %q", actorUsername, movieTitle)
	return s.PushFollowingActivity(ctx, actorID, actorUsername, NotifFollowingLikedReview, &rid, &ref, msg)
}

// PushFollowingAddedWatchlist shortcut สำหรับ watchlist event
func (s *Service) PushFollowingAddedWatchlist(ctx context.Context, actorID uint, actorUsername string, movieID uint, movieTitle string) error {
	mid := movieID
	ref := "movie"
	msg := fmt.Sprintf("%s added %q to their watchlist", actorUsername, movieTitle)
	return s.PushFollowingActivity(ctx, actorID, actorUsername, NotifFollowingAddedWatchlist, &mid, &ref, msg)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (s *Service) toResponse(n Notification) (NotificationResponse, error) {
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
		actor, _, _, _, err := s.userProvider.FindByID(*n.ActorID)
		if err == nil {
			resp.Actor = &ActorSummary{
				ID:          actor.ID,
				Username:    actor.Username,
				DisplayName: actor.DisplayName,
				AvatarURL:   actor.AvatarURL,
			}
		}
	}

	return resp, nil
}

func ptr[T any](v T) *T { return &v }
