package review_module

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrReviewNotFound       = errors.New("review not found")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrForbidden            = errors.New("forbidden")
	ErrAlreadyLiked         = errors.New("already liked")
	ErrNotLiked             = errors.New("not liked")
	ErrAlreadyMarkedHelpful = errors.New("already marked helpful")
	ErrNotMarkedHelpful     = errors.New("not marked helpful")
	ErrInvalidWatchedAt     = errors.New("invalid watched_at")
	ErrInvalidRating        = errors.New("rating must be 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, or 5")
	ErrInvalidMediaType     = errors.New("media_type must be 'movie' or 'tv'")
	ErrInvalidMediaID       = errors.New("invalid media_id")
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// ── Review ────────────────────────────────────────────────────────

func (r *repository) CreateReview(review *Review) error {
	// บันทึกตัวรีวิวเพียวๆ ลงไปก่อน
	return r.db.Create(review).Error
}

func (r *repository) FindReviewsByUser(userID uint) ([]Review, error) {
	var reviews []Review
	// เติม Preload("User") เข้าไปตรงนี้
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *repository) FindReviewsByMedia(mediaID int, mediaType string) ([]Review, error) {
	var reviews []Review
	err := r.db.Preload("User").
		Where("media_id = ? AND media_type = ? AND is_public = true", mediaID, mediaType).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *repository) FindReviewByID(reviewID uint) (*Review, error) {
	var review Review
	// เติม Preload("User") เข้าไปตรงนี้ด้วย เพื่อให้ตอน Update/Delete ทำงานได้ถูกต้อง
	err := r.db.Preload("User").First(&review, reviewID).Error
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

// ── In-app Rating Aggregate ───────────────────────────────────────

// mediaRatingRow เป็น struct รับผลจาก raw query
type mediaRatingRow struct {
	AvgRating   float32
	ReviewCount int
}

// GetMediaRating คืน average rating และจำนวน review ของ media นั้น
// นับเฉพาะ public reviews และ deleted_at IS NULL (soft-delete safe)
func (r *repository) GetMediaRating(mediaID int, mediaType string) (*mediaRatingRow, error) {
	var row mediaRatingRow
	err := r.db.Model(&Review{}).
		Select("COALESCE(AVG(rating), 0) AS avg_rating, COUNT(*) AS review_count").
		Where("media_id = ? AND media_type = ? AND is_public = true AND deleted_at IS NULL",
			mediaID, mediaType).
		Scan(&row).Error
	return &row, err
}

// ── Like ──────────────────────────────────────────────────────────

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

func (r *repository) FindLikedIDs(reviewIDs []uint, userID uint) (map[uint]bool, error) {
	if len(reviewIDs) == 0 {
		return map[uint]bool{}, nil
	}
	var likes []ReviewLike
	err := r.db.Where("review_id IN ? AND user_id = ?", reviewIDs, userID).Find(&likes).Error
	result := make(map[uint]bool)
	for _, l := range likes {
		result[l.ReviewID] = true
	}
	return result, err
}

// ── Helpful ───────────────────────────────────────────────────────
// Pattern เหมือนกับ Like ทุกประการ แยกตารางเพราะคนละความหมาย

func (r *repository) CreateHelpful(reviewID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		vote := ReviewHelpful{ReviewID: reviewID, UserID: userID}
		if err := tx.Create(&vote).Error; err != nil {
			if isDuplicateError(err) {
				return ErrAlreadyMarkedHelpful
			}
			return err
		}
		return tx.Model(&Review{}).Where("id = ?", reviewID).
			UpdateColumn("helpful_count", gorm.Expr("helpful_count + 1")).Error
	})
}

func (r *repository) DeleteHelpful(reviewID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("review_id = ? AND user_id = ?", reviewID, userID).
			Delete(&ReviewHelpful{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotMarkedHelpful
		}
		return tx.Model(&Review{}).Where("id = ?", reviewID).
			UpdateColumn("helpful_count", gorm.Expr("GREATEST(helpful_count - 1, 0)")).Error
	})
}

func (r *repository) IsHelpful(reviewID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ReviewHelpful{}).
		Where("review_id = ? AND user_id = ?", reviewID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) FindHelpfulIDs(reviewIDs []uint, userID uint) (map[uint]bool, error) {
	if len(reviewIDs) == 0 {
		return map[uint]bool{}, nil
	}
	var votes []ReviewHelpful
	err := r.db.Where("review_id IN ? AND user_id = ?", reviewIDs, userID).Find(&votes).Error
	result := make(map[uint]bool)
	for _, v := range votes {
		result[v.ReviewID] = true
	}
	return result, err
}

// ── Comment ───────────────────────────────────────────────────────

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

// ── Helpers ───────────────────────────────────────────────────────

func isDuplicateError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "UNIQUE constraint")
}
