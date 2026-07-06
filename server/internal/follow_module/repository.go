package follow_module

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrAlreadyFollowing = errors.New("already following")
	ErrNotFollowing     = errors.New("not following")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("not found")
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) isPrivate(userID uint) (bool, error) {
	var isPrivate bool
	err := r.db.
		Table("users").
		Select("is_private").
		Where("id = ? AND is_active = true", userID).
		Scan(&isPrivate).Error
	return isPrivate, err
}

func (r *repository) Follow(followerID, followeeID uint) (*UserFollow, error) {
	if followerID == followeeID {
		return nil, ErrForbidden
	}

	var existing UserFollow
	err := r.db.
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		First(&existing).Error
	if err == nil {
		return nil, ErrAlreadyFollowing
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	private, err := r.isPrivate(followeeID)
	if err != nil {
		return nil, err
	}

	status := StatusAccepted
	if private {
		status = StatusPending
	}

	follow := &UserFollow{
		FollowerID: followerID,
		FolloweeID: followeeID,
		Status:     status,
	}
	if err := r.db.Create(follow).Error; err != nil {
		return nil, err
	}
	return follow, nil
}

func (r *repository) Unfollow(followerID, followeeID uint) error {
	result := r.db.
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&UserFollow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFollowing
	}
	return nil
}

func (r *repository) AcceptFollow(followerID, followeeID uint) error {
	result := r.db.
		Model(&UserFollow{}).
		Where("follower_id = ? AND followee_id = ? AND status = ?",
			followerID, followeeID, StatusPending).
		Update("status", StatusAccepted)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) RejectFollow(followerID, followeeID uint) error {
	result := r.db.
		Where("follower_id = ? AND followee_id = ? AND status = ?",
			followerID, followeeID, StatusPending).
		Delete(&UserFollow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

type listRow struct {
	ID          uint
	Username    string
	DisplayName *string
	AvatarURL   *string
	Status      string
}

func (r *repository) GetFollowers(userID uint) ([]listRow, error) {
	var rows []listRow
	err := r.db.
		Table("user_follows uf").
		Select("u.id, u.username, u.display_name, u.avatar_url, uf.status").
		Joins("JOIN users u ON u.id = uf.follower_id AND u.is_active = true").
		Where("uf.followee_id = ? AND uf.status = ?", userID, StatusAccepted).
		Order("uf.created_at DESC").
		Scan(&rows).Error
	return rows, err
}

func (r *repository) GetFollowing(userID uint) ([]listRow, error) {
	var rows []listRow
	err := r.db.
		Table("user_follows uf").
		Select("u.id, u.username, u.display_name, u.avatar_url, uf.status").
		Joins("JOIN users u ON u.id = uf.followee_id AND u.is_active = true").
		Where("uf.follower_id = ? AND uf.status = ?", userID, StatusAccepted).
		Order("uf.created_at DESC").
		Scan(&rows).Error
	return rows, err
}

func (r *repository) GetPendingRequests(userID uint) ([]listRow, error) {
	var rows []listRow
	err := r.db.
		Table("user_follows uf").
		Select("u.id, u.username, u.display_name, u.avatar_url, uf.status").
		Joins("JOIN users u ON u.id = uf.follower_id AND u.is_active = true").
		Where("uf.followee_id = ? AND uf.status = ?", userID, StatusPending).
		Order("uf.created_at DESC").
		Scan(&rows).Error
	return rows, err
}
