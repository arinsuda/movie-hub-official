package user_module

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrForbidden         = errors.New("forbidden")
	ErrOTPNotFound       = errors.New("otp not found")
	ErrOTPExpired        = errors.New("otp expired")
	ErrOTPInvalid        = errors.New("otp invalid")
	ErrOTPMaxAttempts    = errors.New("otp max attempts exceeded")
	ErrEmailAlreadyInUse = errors.New("email already in use")
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func NewPublicRepository(db *gorm.DB) *repository {
	return newRepository(db)
}

// userWithStats เป็น struct รับผลจากการ JOIN users + user_stats view
// ใช้ภายใน repository เท่านั้น ไม่ expose ออก
type userWithStats struct {
	User
	ReviewCount    int `gorm:"column:review_count"`
	FollowerCount  int `gorm:"column:follower_count"`
	FollowingCount int `gorm:"column:following_count"`
}

// FindByID ดึง user พร้อม stats จาก view user_stats
// ข้อมูลจะตรงเสมอ ไม่มี drift จาก denormalized counter
func (r *repository) FindByID(id uint) (*User, int, int, int, error) {
	var row userWithStats

	err := r.db.
		Table("users").
		Select(`
			users.*,
			COALESCE(s.review_count,   0) AS review_count,
			COALESCE(s.follower_count,  0) AS follower_count,
			COALESCE(s.following_count, 0) AS following_count
		`).
		Joins("LEFT JOIN user_stats s ON s.user_id = users.id").
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND users.is_active = true", id).
		First(&row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, 0, 0, ErrUserNotFound
	}
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// โหลด Role แยก (Preload ไม่ทำงานกับ Table() โดยตรง)
	if err := r.db.First(&row.User.Role, row.User.RoleID).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return &row.User, row.ReviewCount, row.FollowerCount, row.FollowingCount, nil
}

func (r *repository) UpdateProfile(id uint, updates map[string]any) error {
	result := r.db.Model(&User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) DeleteUser(id uint) error {
	result := r.db.Model(&User{}).Where("id = ?", id).Update("is_active", false)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

// ── Email Change OTP ──────────────────────────────────────────────

// UpsertEmailChangeRequest สร้างหรือแทนที่ pending request ของ user
// ใช้ upsert เพื่อให้ request ใหม่ยกเลิกของเก่าอัตโนมัติ
func (r *repository) UpsertEmailChangeRequest(req *EmailChangeRequest) error {
	return r.db.
		Where(EmailChangeRequest{UserID: req.UserID}).
		Assign(EmailChangeRequest{
			OTPHash:      req.OTPHash,
			ExpiresAt:    req.ExpiresAt,
			AttemptCount: 0,
		}).
		FirstOrCreate(req).Error
}

func (r *repository) FindEmailChangeRequest(userID uint) (*EmailChangeRequest, error) {
	var req EmailChangeRequest
	err := r.db.Where("user_id = ?", userID).First(&req).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOTPNotFound
	}
	return &req, err
}

func (r *repository) IncrementOTPAttempt(id uint) error {
	return r.db.Model(&EmailChangeRequest{}).
		Where("id = ?", id).
		UpdateColumn("attempt_count", gorm.Expr("attempt_count + 1")).Error
}

func (r *repository) DeleteEmailChangeRequest(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&EmailChangeRequest{}).Error
}

// IsEmailTaken ตรวจว่า email ถูกใช้งานแล้วหรือยัง (ยกเว้น excludeUserID)
func (r *repository) IsEmailTaken(email string, excludeUserID uint) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).
		Where("email = ? AND id != ? AND is_active = true", email, excludeUserID).
		Count(&count).Error
	return count > 0, err
}
