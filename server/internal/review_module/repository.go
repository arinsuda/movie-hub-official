package review_module

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrReviewNotFound   = errors.New("review not found")
	ErrCommentNotFound  = errors.New("comment not found")
	ErrForbidden        = errors.New("forbidden")
	ErrAlreadyLiked     = errors.New("already liked")
	ErrNotLiked         = errors.New("not liked")
	ErrInvalidWatchedAt = errors.New("invalid watched_at")
	ErrInvalidRating    = errors.New("rating must be between 0 and 10")
	ErrInvalidMediaType = errors.New("media_type must be 'movie' or 'tv'")
	ErrInvalidMediaID   = errors.New("invalid media_id")
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// Review
func (r *repository) CreateReview(review *Review) error {
	return r.db.Create(review).Error
}

func (r *repository) FindReviewsByUser(userID uint) ([]Review, error) {
	var reviews []Review
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *repository) FindReviewsByMedia(mediaID int, mediaType string) ([]Review, error) {
	var reviews []Review
	err := r.db.Where("media_id = ? AND media_type = ? AND is_public = true", mediaID, mediaType).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *repository) FindReviewByID(reviewID uint) (*Review, error) {
	var review Review
	err := r.db.First(&review, reviewID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrReviewNotFound
	}
	return &review, err
}

func (r *repository) UpdateReview(reviewID uint, updates map[string]any) error {
	result := r.db.Model(&Review{}).Where("id = ?", reviewID).Updates(updates)
	if result.RowsAffected == 0 {
		return ErrReviewNotFound
	}
	return result.Error
}

func (r *repository) DeleteReview(reviewID uint) error {
	result := r.db.Delete(&Review{}, reviewID)
	if result.RowsAffected == 0 {
		return ErrReviewNotFound
	}
	return result.Error
}

// Like
func (r *repository) CreateLike(reviewID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		like := ReviewLike{ReviewID: reviewID, UserID: userID}
		if err := tx.Create(&like).Error; err != nil {
			if isDuplicateError(err) {
				return ErrAlreadyLiked
			}
			return err
		}
		return tx.Model(&Review{}).Where("id = ?", reviewID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
	})
}

func (r *repository) DeleteLike(reviewID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("review_id = ? AND user_id = ?", reviewID, userID).
			Delete(&ReviewLike{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotLiked
		}
		return tx.Model(&Review{}).Where("id = ?", reviewID).
			UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
	})
}

func (r *repository) IsLiked(reviewID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ReviewLike{}).
		Where("review_id = ? AND user_id = ?", reviewID, userID).
		Count(&count).Error
	return count > 0, err
}

// Comment
func (r *repository) CreateComment(comment *ReviewComment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		return tx.Model(&Review{}).Where("id = ?", comment.ReviewID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
	})
}

func (r *repository) FindCommentsByReview(reviewID uint) ([]ReviewComment, error) {
	var comments []ReviewComment
	err := r.db.Where("review_id = ?", reviewID).
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

func (r *repository) FindCommentByID(commentID uint) (*ReviewComment, error) {
	var comment ReviewComment
	err := r.db.First(&comment, commentID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCommentNotFound
	}
	return &comment, err
}

func (r *repository) UpdateComment(commentID uint, body string) error {
	result := r.db.Model(&ReviewComment{}).Where("id = ?", commentID).Update("body", body)
	if result.RowsAffected == 0 {
		return ErrCommentNotFound
	}
	return result.Error
}

func (r *repository) DeleteComment(commentID uint, reviewID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND review_id = ?", commentID, reviewID).Delete(&ReviewComment{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrCommentNotFound
		}
		return tx.Model(&Review{}).Where("id = ?", reviewID).
			UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error
	})
}

func isDuplicateError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "UNIQUE constraint")
}

func (r *repository) FindLikedIDs(reviewIDs []uint, userID uint) (map[uint]bool, error) {
	var likes []ReviewLike
	err := r.db.Where("review_id IN ? AND user_id = ?", reviewIDs, userID).Find(&likes).Error
	result := make(map[uint]bool)
	for _, l := range likes {
		result[l.ReviewID] = true
	}
	return result, err
}
