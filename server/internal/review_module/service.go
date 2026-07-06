package review_module

import (
	"context"
	"math"
	"time"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	tmdbmodule "github.com/arinsuda/movie-hub/internal/tmdb_module"
	users "github.com/arinsuda/movie-hub/internal/user_module"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"gorm.io/gorm"
)

type Service struct {
	repo       *repository
	db         *gorm.DB
	minio      *storage.MinIOClient
	expPort    stats.ExpAdder
	achieveSvc achievementsmodule.Service
	notifSvc   *notification_module.Service
}

func NewService(
	db *gorm.DB,
	mc *storage.MinIOClient,
	exp stats.ExpAdder,
	achieve achievementsmodule.Service,
	notif *notification_module.Service,
) *Service {
	return &Service{
		repo:       newRepository(db),
		db:         db,
		minio:      mc,
		expPort:    exp,
		achieveSvc: achieve,
		notifSvc:   notif,
	}
}

// ── Review ────────────────────────────────────────────────────────

func (s *Service) CreateReview(userID uint, req CreateReviewRequest) (*ReviewResponse, error) {
	if err := validateReviewRequest(req); err != nil {
		return nil, err
	}

	review := &Review{
		UserID:    userID,
		MediaID:   req.MediaID,
		MediaType: req.MediaType,
		Rating:    req.Rating,
		Body:      req.Body,
		IsPublic:  req.IsPublic,
	}
	if req.WatchedAt != nil {
		t, err := time.Parse("2006-01-02", *req.WatchedAt)
		if err != nil {
			return nil, ErrInvalidWatchedAt
		}
		review.WatchedAt = &t
	}

	if err := s.repo.CreateReview(review); err != nil {
		return nil, err
	}

	// ── EXP ───────────────────────────────────────────────────────
	_ = s.expPort.AddExperience(userID, stats.ExpPerReview)

	// ── Achievement: review_count ─────────────────────────────────
	var reviewCount int64
	s.db.Model(&Review{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&reviewCount)
	_, _ = s.achieveSvc.Track(userID, "review_count", int(reviewCount))

	// ── Achievement: rating_one_star_count / rating_five_star_count ──
	if req.Rating <= 1.0 {
		var oneStarCount int64
		s.db.Model(&Review{}).
			Where("user_id = ? AND rating <= 1.0 AND deleted_at IS NULL", userID).
			Count(&oneStarCount)
		_, _ = s.achieveSvc.Track(userID, "rating_one_star_count", int(oneStarCount))
	}
	if req.Rating == 5.0 {
		var fiveStarCount int64
		s.db.Model(&Review{}).
			Where("user_id = ? AND rating = 5.0 AND deleted_at IS NULL", userID).
			Count(&fiveStarCount)
		_, _ = s.achieveSvc.Track(userID, "rating_five_star_count", int(fiveStarCount))
	}
	if req.Rating < 3.0 {
		var lowRatingCount int64
		s.db.Model(&Review{}).
			Where("user_id = ? AND rating < 3.0 AND deleted_at IS NULL", userID).
			Count(&lowRatingCount)
		_, _ = s.achieveSvc.Track(userID, "low_rating_count", int(lowRatingCount))
	}

	// ── Notification: fan-out ให้ followers ──────────────────────
	if req.IsPublic && s.notifSvc != nil {
		if actor, err := s.getUserSummary(userID); err == nil {
			title, _ := fetchMediaSummary(req.MediaID, req.MediaType)
			_ = s.notifSvc.PushFollowingReviewed(
				context.Background(),
				userID,
				actor.Username,
				review.ID,
				title,
			)
		}
	}

	inserted, err := s.repo.FindReviewByID(review.ID)
	if err != nil {
		return nil, err
	}
	// review ที่เพิ่งสร้าง ยังไม่มีใคร like/helpful แน่นอน (รวมถึงตัวเอง)
	return toReviewResponse(inserted, false, false, s.minio), nil
}

// GetUserReviews คืนรีวิวของ userID ตาม filter ที่กำหนด (visibility + date range)
// ถ้า requesterID ไม่ใช่เจ้าของ (userID) จะบังคับเห็นเฉพาะ public เท่านั้น
// ไม่ว่า filter.Visibility จะขอ "private" หรือ "all" มาก็ตาม
func (s *Service) GetUserReviews(userID, requesterID uint, filter ReviewFilter) ([]ReviewResponse, error) {
	if userID != requesterID {
		filter.Visibility = "public"
	}

	reviews, err := s.repo.FindReviewsByUser(userID, filter)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(reviews))
	for i, r := range reviews {
		ids[i] = r.ID
	}

	likedMap, _ := s.repo.FindLikedIDs(ids, requesterID)
	helpfulMap, _ := s.repo.FindHelpfulIDs(ids, requesterID)

	responses := make([]ReviewResponse, len(reviews))
	for i, r := range reviews {
		responses[i] = *toReviewResponse(&r, likedMap[r.ID], helpfulMap[r.ID], s.minio)
	}
	return responses, nil
}

func (s *Service) GetMediaReviews(mediaID int, mediaType string, requesterID uint) ([]ReviewResponse, error) {
	reviews, err := s.repo.FindReviewsByMedia(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(reviews))
	for i, r := range reviews {
		ids[i] = r.ID
	}
	likedMap, _ := s.repo.FindLikedIDs(ids, requesterID)
	helpfulMap, _ := s.repo.FindHelpfulIDs(ids, requesterID)

	responses := make([]ReviewResponse, len(reviews))
	for i, r := range reviews {
		responses[i] = *toReviewResponse(&r, likedMap[r.ID], helpfulMap[r.ID], s.minio)
	}
	return responses, nil
}

func (s *Service) UpdateReview(reviewID, requesterID uint, req UpdateReviewRequest) (*ReviewResponse, error) {
	review, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return nil, err
	}
	if review.UserID != requesterID {
		return nil, ErrForbidden
	}

	updates := buildUpdateMap(req)
	if len(updates) == 0 {
		liked, _ := s.repo.IsLiked(reviewID, requesterID)
		helpful, _ := s.repo.IsHelpful(reviewID, requesterID)
		return toReviewResponse(review, liked, helpful, s.minio), nil
	}

	if err := s.repo.UpdateReview(reviewID, updates); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return nil, err
	}
	liked, _ := s.repo.IsLiked(reviewID, requesterID)
	helpful, _ := s.repo.IsHelpful(reviewID, requesterID)
	return toReviewResponse(updated, liked, helpful, s.minio), nil
}

func (s *Service) DeleteReview(reviewID, requesterID uint) error {
	review, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return err
	}
	if review.UserID != requesterID {
		return ErrForbidden
	}

	if err := s.repo.DeleteReview(reviewID); err != nil {
		return err
	}

	_ = s.expPort.AddExperience(requesterID, -stats.ExpPerReview)
	return nil
}

// ── In-app Rating ─────────────────────────────────────────────────

func (s *Service) GetMediaRating(mediaID int, mediaType string) (*RatingResponse, error) {
	if mediaType != "movie" && mediaType != "tv" {
		return nil, ErrInvalidMediaType
	}
	if mediaID <= 0 {
		return nil, ErrInvalidMediaID
	}

	row, err := s.repo.GetMediaRating(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	hasRating := row.ReviewCount > 0
	avg := float32(0)
	if hasRating {
		avg = float32(math.Round(float64(row.AvgRating)*10) / 10)
	}

	return &RatingResponse{
		MediaID:       mediaID,
		MediaType:     mediaType,
		AverageRating: avg,
		ReviewCount:   row.ReviewCount,
		HasRating:     hasRating,
	}, nil
}

// ── Like ──────────────────────────────────────────────────────────

func (s *Service) LikeReview(reviewID, requesterID uint) error {
	review, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return err
	}
	if !review.IsPublic && review.UserID != requesterID {
		return ErrForbidden
	}

	if err := s.repo.CreateLike(reviewID, requesterID); err != nil {
		return err
	}

	// ── EXP ──────────────────────────────────────────────────────
	_ = s.expPort.AddExperience(review.UserID, stats.ExpPerLike)

	// ── Achievement: review_like_given_count (ฝั่งคนกด like) ────
	var likeGivenCount int64
	s.db.Model(&ReviewLike{}).Where("user_id = ?", requesterID).Count(&likeGivenCount)
	_, _ = s.achieveSvc.Track(requesterID, "review_like_given_count", int(likeGivenCount))

	// ── Achievement: review_like_received_count (ฝั่งเจ้าของ review) ──
	var likeReceivedCount int64
	s.db.Model(&ReviewLike{}).
		Joins("JOIN reviews ON reviews.id = review_likes.review_id").
		Where("reviews.user_id = ? AND reviews.deleted_at IS NULL", review.UserID).
		Count(&likeReceivedCount)
	_, _ = s.achieveSvc.Track(review.UserID, "review_like_received_count", int(likeReceivedCount))

	// ── Notification: แจ้งเจ้าของ review (ถ้าไม่ใช่คนเดียวกัน) ──
	if s.notifSvc != nil && review.UserID != requesterID {
		if actor, err := s.getUserSummary(requesterID); err == nil {
			title, _ := fetchMediaSummary(review.MediaID, review.MediaType)
			_ = s.notifSvc.PushFollowingLikedReview(
				context.Background(),
				requesterID,
				actor.Username,
				reviewID,
				title,
			)
		}
	}

	return nil
}

func (s *Service) UnlikeReview(reviewID, requesterID uint) error {
	review, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteLike(reviewID, requesterID); err != nil {
		return err
	}

	_ = s.expPort.AddExperience(review.UserID, -stats.ExpPerLike)
	return nil
}

// ── Helpful ───────────────────────────────────────────────────────
// แนวคิดแบบ Stack Overflow: โหวตว่ารีวิวนี้ "มีประโยชน์" มากแค่ไหน
// ไม่ผูกกับความชอบส่วนตัว (like) แยกกันชัดเจน

func (s *Service) MarkHelpful(reviewID, requesterID uint) error {
	review, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return err
	}
	if !review.IsPublic && review.UserID != requesterID {
		return ErrForbidden
	}
	// เจ้าของรีวิวโหวต helpful รีวิวตัวเองไม่ได้ (กันการปั่นตัวเลข)
	if review.UserID == requesterID {
		return ErrForbidden
	}

	if err := s.repo.CreateHelpful(reviewID, requesterID); err != nil {
		return err
	}

	// ── Achievement: review_helpful_received_count (ฝั่งเจ้าของ review) ──
	var helpfulReceivedCount int64
	s.db.Model(&ReviewHelpful{}).
		Joins("JOIN reviews ON reviews.id = review_helpfuls.review_id").
		Where("reviews.user_id = ? AND reviews.deleted_at IS NULL", review.UserID).
		Count(&helpfulReceivedCount)
	_, _ = s.achieveSvc.Track(review.UserID, "review_helpful_received_count", int(helpfulReceivedCount))

	// ── Notification: แจ้งเจ้าของ review (ถ้าไม่ใช่คนเดียวกัน) ──
	if s.notifSvc != nil {
		if actor, err := s.getUserSummary(requesterID); err == nil {
			title, _ := fetchMediaSummary(review.MediaID, review.MediaType)
			_ = s.notifSvc.PushFollowingLikedReview(
				context.Background(),
				requesterID,
				actor.Username,
				reviewID,
				title,
			)
		}
	}

	return nil
}

func (s *Service) UnmarkHelpful(reviewID, requesterID uint) error {
	if _, err := s.repo.FindReviewByID(reviewID); err != nil {
		return err
	}
	return s.repo.DeleteHelpful(reviewID, requesterID)
}

// ── Comment ───────────────────────────────────────────────────────

func (s *Service) CreateComment(reviewID, requesterID uint, req CreateCommentRequest) (*CommentResponse, error) {
	review, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return nil, err
	}
	if !review.IsPublic && review.UserID != requesterID {
		return nil, ErrForbidden
	}

	comment := &ReviewComment{ReviewID: reviewID, UserID: requesterID, Body: req.Body}
	if err := s.repo.CreateComment(comment); err != nil {
		return nil, err
	}
	return toCommentResponse(comment), nil
}

func (s *Service) GetComments(reviewID, requesterID uint) ([]CommentResponse, error) {
	review, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return nil, err
	}
	if !review.IsPublic && review.UserID != requesterID {
		return nil, ErrForbidden
	}

	comments, err := s.repo.FindCommentsByReview(reviewID)
	if err != nil {
		return nil, err
	}

	responses := make([]CommentResponse, len(comments))
	for i, c := range comments {
		responses[i] = *toCommentResponse(&c)
	}
	return responses, nil
}

func (s *Service) UpdateComment(commentID, requesterID uint, req UpdateCommentRequest) (*CommentResponse, error) {
	comment, err := s.repo.FindCommentByID(commentID)
	if err != nil {
		return nil, err
	}
	if comment.UserID != requesterID {
		return nil, ErrForbidden
	}
	if err := s.repo.UpdateComment(commentID, req.Body); err != nil {
		return nil, err
	}
	comment.Body = req.Body
	return toCommentResponse(comment), nil
}

func (s *Service) DeleteComment(commentID, reviewID, requesterID uint) error {
	comment, err := s.repo.FindCommentByID(commentID)
	if err != nil {
		return err
	}
	if comment.UserID != requesterID {
		return ErrForbidden
	}
	return s.repo.DeleteComment(commentID, reviewID)
}

// ── Helpers ───────────────────────────────────────────────────────

func (s *Service) getUserSummary(userID uint) (*users.User, error) {
	var u users.User
	err := s.db.First(&u, userID).Error
	return &u, err
}

func validateReviewRequest(req CreateReviewRequest) error {
	if req.Rating < 0.5 || req.Rating > 5 || math.Mod(float64(req.Rating)*2, 1) != 0 {
		return ErrInvalidRating
	}
	if req.MediaType != "movie" && req.MediaType != "tv" {
		return ErrInvalidMediaType
	}
	if req.MediaID <= 0 {
		return ErrInvalidMediaID
	}
	return nil
}

func buildUpdateMap(req UpdateReviewRequest) map[string]any {
	updates := map[string]any{}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Body != nil {
		updates["body"] = *req.Body
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.WatchedAt != nil {
		t, err := time.Parse("2006-01-02", *req.WatchedAt)
		if err == nil {
			updates["watched_at"] = t
		}
	}
	return updates
}

func toReviewResponse(r *Review, isLiked, isHelpful bool, minio *storage.MinIOClient) *ReviewResponse {
	if r == nil {
		return nil
	}

	title, posterURL := fetchMediaSummary(r.MediaID, r.MediaType)

	return &ReviewResponse{
		ID:   r.ID,
		User: toUserSummaryResponse(&r.User, minio),
		Media: tmdbmodule.Media{
			ID:        r.MediaID,
			MediaType: r.MediaType,
			Title:     title,
			PosterURL: posterURL,
		},
		Rating:         r.Rating,
		Body:           r.Body,
		IsPublic:       r.IsPublic,
		WatchedAt:      r.WatchedAt,
		LikeCount:      r.LikeCount,
		CommentCount:   r.CommentCount,
		IsLiked:        isLiked,
		HelpfulCount:   r.HelpfulCount,
		IsHelpfulVoted: isHelpful,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func toUserSummaryResponse(u *users.User, minio *storage.MinIOClient) users.UserSummaryResponse {
	if u == nil || u.ID == 0 {
		return users.UserSummaryResponse{Username: "Unknown User"}
	}
	res := users.UserSummaryResponse{ID: u.ID, Username: u.Username}
	if u.DisplayName != nil && *u.DisplayName != "" {
		res.DisplayName = u.DisplayName
	}
	if u.AvatarURL != nil && *u.AvatarURL != "" {
		if minio != nil {
			if presigned, err := minio.PresignURL(context.Background(), *u.AvatarURL); err == nil {
				res.AvatarURL = &presigned
			} else {
				res.AvatarURL = u.AvatarURL
			}
		} else {
			res.AvatarURL = u.AvatarURL
		}
	}
	return res
}

func fetchMediaSummary(mediaID int, mediaType string) (title, posterURL string) {
	if mediaType == "movie" {
		movie, err := tmdbmodule.GetMovieByID(mediaID)
		if err != nil {
			return "", ""
		}
		return movie.Title, tmdbmodule.ImageURL(movie.PosterPath)
	}

	series, err := tmdbmodule.GetSeriesByID(mediaID)
	if err != nil {
		return "", ""
	}
	return series.Name, tmdbmodule.ImageURL(series.PosterPath)
}

func toCommentResponse(c *ReviewComment) *CommentResponse {
	return &CommentResponse{
		ID:        c.ID,
		ReviewID:  c.ReviewID,
		UserID:    c.UserID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
