package review_module

import (
	"math"
	"time"

	statsmodule "github.com/arinsuda/movie-hub/internal/stats_module"
	"gorm.io/gorm"
)

type Service struct {
	repo  *repository
	stats *statsmodule.Service
}

func NewService(db *gorm.DB, stats *statsmodule.Service) *Service {
	return &Service{
		repo:  newRepository(db),
		stats: stats,
	}
}

// ── Review ────────────────────────────────────────────────────────

func (s *Service) CreateReview(userID uint, req CreateReviewRequest) (*ReviewResponse, error) {
	if req.Rating < 0 || req.Rating > 10 {
		return nil, ErrInvalidRating
	}
	if req.MediaType != "movie" && req.MediaType != "tv" {
		return nil, ErrInvalidMediaType
	}
	if req.MediaID <= 0 {
		return nil, ErrInvalidMediaID
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

	// Sync review_count in stats (best-effort — don't fail the main operation)
	_ = s.stats.IncrementReviewCount(req.MediaID, req.MediaType, 1)

	return toReviewResponse(review, false), nil
}

func (s *Service) GetUserReviews(userID uint, requesterID uint) ([]ReviewResponse, error) {
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
		if err != nil {
			return nil, ErrInvalidWatchedAt
		}
		updates["watched_at"] = t
	}

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

	// Sync review_count in stats (best-effort)
	_ = s.stats.IncrementReviewCount(review.MediaID, review.MediaType, -1)

	return nil
}

// ── In-app Rating Aggregate ───────────────────────────────────────

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
	return s.repo.CreateLike(reviewID, requesterID)
}

func (s *Service) UnlikeReview(reviewID, requesterID uint) error {
	if _, err := s.repo.FindReviewByID(reviewID); err != nil {
		return err
	}
	return s.repo.DeleteLike(reviewID, requesterID)
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

	comment := &ReviewComment{
		ReviewID: reviewID,
		UserID:   requesterID,
		Body:     req.Body,
	}
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

func toReviewResponse(r *Review, isLiked bool) *ReviewResponse {
	return &ReviewResponse{
		ID:           r.ID,
		UserID:       r.UserID,
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
