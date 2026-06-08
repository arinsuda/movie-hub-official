package review_module

import (
	"math"
	"time"

	users "github.com/arinsuda/movie-hub/internal/user_module"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"gorm.io/gorm"
)

type Service struct {
	repo    *repository
	expPort stats.ExpAdder // port — no direct import of the concrete Service
}

func NewService(db *gorm.DB, exp stats.ExpAdder) *Service {
	return &Service{repo: newRepository(db), expPort: exp}
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

	// Award EXP — best-effort, do not fail the request on error.
	_ = s.expPort.AddExperience(userID, stats.ExpPerReview)

	inserted, err := s.repo.FindReviewByID(review.ID)
	if err != nil {
		return nil, err
	}
	return toReviewResponse(inserted, false), nil
}

func (s *Service) GetUserReviews(userID, requesterID uint) ([]ReviewResponse, error) {
	reviews, err := s.repo.FindReviewsByUser(userID)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(reviews))
	for _, r := range reviews {
		if userID != requesterID && !r.IsPublic {
			continue
		}
		ids = append(ids, r.ID)
	}
	likedMap, _ := s.repo.FindLikedIDs(ids, requesterID)

	responses := make([]ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		if userID != requesterID && !r.IsPublic {
			continue
		}
		responses = append(responses, *toReviewResponse(&r, likedMap[r.ID]))
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

	responses := make([]ReviewResponse, len(reviews))
	for i, r := range reviews {
		responses[i] = *toReviewResponse(&r, likedMap[r.ID])
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
		return toReviewResponse(review, liked), nil
	}

	if err := s.repo.UpdateReview(reviewID, updates); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindReviewByID(reviewID)
	if err != nil {
		return nil, err
	}
	liked, _ := s.repo.IsLiked(reviewID, requesterID)
	return toReviewResponse(updated, liked), nil
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

	// Revoke EXP — best-effort.
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

	// Award EXP to the review author, not the liker.
	_ = s.expPort.AddExperience(review.UserID, stats.ExpPerLike)
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

	// Revoke EXP from the review author.
	_ = s.expPort.AddExperience(review.UserID, -stats.ExpPerLike)
	return nil
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

func toReviewResponse(r *Review, isLiked bool) *ReviewResponse {
	if r == nil {
		return nil
	}
	return &ReviewResponse{
		ID:           r.ID,
		User:         toUserSummaryResponse(&r.User),
		MediaID:      r.MediaID,
		MediaType:    r.MediaType,
		Rating:       r.Rating,
		Body:         r.Body,
		IsPublic:     r.IsPublic,
		WatchedAt:    r.WatchedAt,
		LikeCount:    r.LikeCount,
		CommentCount: r.CommentCount,
		IsLiked:      isLiked,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
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

func toUserSummaryResponse(u *users.User) users.UserSummaryResponse {
	if u == nil || u.ID == 0 {
		return users.UserSummaryResponse{Username: "Unknown User"}
	}
	res := users.UserSummaryResponse{ID: u.ID, Username: u.Username}
	if u.DisplayName != nil && *u.DisplayName != "" {
		res.DisplayName = u.DisplayName
	}
	if u.AvatarURL != nil && *u.AvatarURL != "" {
		res.AvatarURL = u.AvatarURL
	}
	return res
}
