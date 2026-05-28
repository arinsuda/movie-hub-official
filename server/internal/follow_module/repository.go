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

// isPrivate เช็คว่า user นั้น private ไหม
// return error ถ้า user ไม่มีอยู่หรือ inactive
func (r *repository) isPrivate(userID uint) (bool, error) {
	var isPrivate bool
	err := r.db.
		Table("users").
		Select("is_private").
		Where("id = ? AND is_active = true", userID).
		Scan(&isPrivate).Error
	return isPrivate, err
}

// Follow สร้าง follow record
// status = accepted ถ้า followee เป็น public
// status = pending  ถ้า followee เป็น private
func (r *repository) Follow(followerID, followeeID uint) (*UserFollow, error) {
	if followerID == followeeID {
		return nil, ErrForbidden
	}

	// เช็ค duplicate ก่อน (uniqueIndex จะ catch ได้อยู่แล้ว แต่ error message จะไม่ friendly)
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

// Unfollow ลบ record ไม่ว่า status จะเป็นอะไร (cancel pending ได้ด้วย)
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

// AcceptFollow เปลี่ยน pending → accepted
// followeeID = คนที่กด accept (เจ้าของ account)
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

// RejectFollow ลบ pending request
// followeeID = คนที่กด reject
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

// listRow ใช้ scan ผล JOIN กับ users table
type listRow struct {
	ID          uint
	Username    string
	DisplayName *string
	AvatarURL   *string
	Status      string
}

// GetFollowers คืนรายชื่อคนที่ follow userID (accepted เท่านั้น)
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

// GetFollowing คืนรายชื่อคนที่ userID follow (accepted เท่านั้น)
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

// GetPendingRequests คืน pending requests ที่รอ userID approve
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
