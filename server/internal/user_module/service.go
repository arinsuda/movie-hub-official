package user_module

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	"github.com/arinsuda/movie-hub/internal/shared/storage"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo     *repository
	minio    *storage.MinIOClient
	statsSvc *stats.Service
	mailer   Mailer // เพิ่มตรงนี้
}

// แก้ NewService
func NewService(db *gorm.DB, mc *storage.MinIOClient, statsSvc *stats.Service, mailer Mailer) *Service {
	return &Service{
		repo:     newRepository(db),
		minio:    mc,
		statsSvc: statsSvc,
		mailer:   mailer,
	}
}

// ── GetProfile ────────────────────────────────────────────────────

func (s *Service) GetProfile(targetUserID, requesterID uint) (*UserProfileResponse, error) {
	user, reviewCount, followerCount, followingCount, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	level := s.fetchLevel(targetUserID)
	profile := toProfileResponse(user, reviewCount, followerCount, followingCount, level)

	if user.IsPrivate && requesterID != targetUserID {
		profile.Bio = nil
		profile.FavoriteGenres = nil
		profile.Gender = ""
		profile.DateOfBirth = nil
	}

	if profile.AvatarURL != nil && strings.HasPrefix(*profile.AvatarURL, "avatars/") {
		if url, err := s.minio.PresignURL(context.Background(), *profile.AvatarURL); err == nil {
			profile.AvatarURL = &url
		}
	}

	return profile, nil
}

// ── UpdateProfile ─────────────────────────────────────────────────

func (s *Service) UpdateProfile(
	targetUserID, requesterID uint,
	req UpdateProfileRequest,
	avatarFile multipart.File,
	avatarHeader *multipart.FileHeader,
) (*UserProfileResponse, error) {
	if targetUserID != requesterID {
		return nil, ErrForbidden
	}

	updates := map[string]any{}
	if req.DisplayName != nil {
		updates["display_name"] = req.DisplayName
	}
	if req.Bio != nil {
		updates["bio"] = req.Bio
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = req.DateOfBirth
	}
	if req.Gender != "" {
		if !isValidGender(req.Gender) {
			return nil, errors.New("invalid gender")
		}
		updates["gender"] = req.Gender
	}
	if req.GenderOther != nil {
		updates["gender_other"] = req.GenderOther
	}
	if req.FavoriteGenres != nil {
		updates["favorite_genres"] = req.FavoriteGenres
	}
	if req.IsPrivate != nil {
		updates["is_private"] = req.IsPrivate
	}

	if avatarFile != nil && avatarHeader != nil {
		oldKey := s.currentAvatarKey(targetUserID)
		newKey, err := s.minio.UploadAvatar(context.Background(), targetUserID, avatarFile, avatarHeader)
		if err != nil {
			return nil, err
		}
		updates["avatar_url"] = newKey
		if oldKey != "" {
			_ = s.minio.DeleteObject(context.Background(), oldKey)
		}
	}
	log.Printf("DEBUG updates map: %+v", updates)
	if len(updates) > 0 {
		if err := s.repo.UpdateProfile(targetUserID, updates); err != nil {
			return nil, err
		}
	}

	return s.GetProfile(targetUserID, requesterID)
}

// ── DeleteUser ────────────────────────────────────────────────────

func (s *Service) DeleteUser(targetUserID, requesterID uint, requesterRole string) error {
	if targetUserID != requesterID && requesterRole != "admin" {
		return ErrForbidden
	}
	return s.repo.DeleteUser(targetUserID)
}

// ── UpdateFavoriteGenres ──────────────────────────────────────────

func (s *Service) UpdateFavoriteGenres(targetUserID, requesterID uint, genres []int) (*UserProfileResponse, error) {
	if targetUserID != requesterID {
		return nil, ErrForbidden
	}

	var stored *string
	if len(genres) > 0 {
		b, err := json.Marshal(genres)
		if err != nil {
			return nil, err
		}
		str := string(b)
		stored = &str
	}

	if err := s.repo.UpdateProfile(targetUserID, map[string]any{"favorite_genres": stored}); err != nil {
		return nil, err
	}
	return s.GetProfile(targetUserID, requesterID)
}

// ── RequestEmailChange ────────────────────────────────────────────

func (s *Service) RequestEmailChange(targetUserID, requesterID uint, newEmail string) error {
	if targetUserID != requesterID {
		return ErrForbidden
	}

	// ตรวจ email format เบื้องต้น
	newEmail = strings.TrimSpace(strings.ToLower(newEmail))
	if !isValidEmail(newEmail) {
		return errors.New("invalid email format")
	}

	// โหลด user ปัจจุบันเพื่อดึง email จริง
	user, _, _, _, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return err
	}

	// ห้ามเปลี่ยนเป็น email เดิม
	if strings.EqualFold(user.Email, newEmail) {
		return errors.New("new email must be different from current email")
	}

	// ตรวจว่า email ใหม่ถูกใช้งานแล้วหรือยัง
	taken, err := s.repo.IsEmailTaken(newEmail, targetUserID)
	if err != nil {
		return err
	}
	if taken {
		return ErrEmailAlreadyInUse
	}

	// สร้าง OTP 6 หลัก
	otp, err := generateOTP()
	if err != nil {
		return err
	}

	// Hash OTP ก่อนเก็บ (ใช้ bcrypt)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Upsert pending request (แทนของเก่าถ้ามี)
	changeReq := &EmailChangeRequest{
		UserID:    targetUserID,
		NewEmail:  newEmail,
		OTPHash:   string(otpHash),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.repo.UpsertEmailChangeRequest(changeReq); err != nil {
		return err
	}

	// ส่ง OTP ไปยัง email ปัจจุบัน (ไม่ใช่ email ใหม่!)
	if err := s.mailer.SendOTP(user.Email, otp); err != nil {
		// ถ้าส่งเมลไม่ได้ ลบ record ทิ้ง ไม่ให้ค้างไว้
		_ = s.repo.DeleteEmailChangeRequest(targetUserID)
		return errors.New("failed to send OTP email")
	}

	return nil
}

// ── VerifyEmailChange ─────────────────────────────────────────────

const maxOTPAttempts = 5

func (s *Service) VerifyEmailChange(targetUserID, requesterID uint, otp string) (*UserProfileResponse, error) {
	if targetUserID != requesterID {
		return nil, ErrForbidden
	}

	// ดึง pending request
	req, err := s.repo.FindEmailChangeRequest(targetUserID)
	if err != nil {
		return nil, err // ErrOTPNotFound
	}

	// ตรวจหมดเวลา
	if req.IsExpired() {
		_ = s.repo.DeleteEmailChangeRequest(targetUserID)
		return nil, ErrOTPExpired
	}

	// ตรวจจำนวนครั้งที่กรอกผิด
	if req.AttemptCount >= maxOTPAttempts {
		_ = s.repo.DeleteEmailChangeRequest(targetUserID)
		return nil, ErrOTPMaxAttempts
	}

	// ตรวจ OTP
	if err := bcrypt.CompareHashAndPassword([]byte(req.OTPHash), []byte(otp)); err != nil {
		// นับความผิดพลาด
		_ = s.repo.IncrementOTPAttempt(req.ID)
		return nil, ErrOTPInvalid
	}

	// OTP ถูก → ตรวจ email ใหม่ว่าถูกแย่งไปแล้วหรือยัง (race condition guard)
	taken, err := s.repo.IsEmailTaken(req.NewEmail, targetUserID)
	if err != nil {
		return nil, err
	}
	if taken {
		_ = s.repo.DeleteEmailChangeRequest(targetUserID)
		return nil, ErrEmailAlreadyInUse
	}

	// อัปเดต email จริงในตาราง users
	if err := s.repo.UpdateProfile(targetUserID, map[string]any{
		"email":             req.NewEmail,
		"verified_email_at": time.Now(), // reset verification
	}); err != nil {
		return nil, err
	}

	// ลบ pending request
	_ = s.repo.DeleteEmailChangeRequest(targetUserID)

	return s.GetProfile(targetUserID, requesterID)
}

// ── Helpers ───────────────────────────────────────────────────────

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	// zero-pad เป็น 6 หลักเสมอ เช่น 000042
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func isValidEmail(email string) bool {
	// regex เบื้องต้น ครอบคลุม 99% ของ use case
	const pattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ── Helpers ───────────────────────────────────────────────────────

// fetchLevel reads the user's current level from user_stats_module.
// Returns 1 on any error so the profile response is never broken.
func (s *Service) fetchLevel(userID uint) int {
	if s.statsSvc == nil {
		return 1
	}
	st, err := s.statsSvc.GetUserStats(userID)
	if err != nil || st == nil {
		return 1
	}
	return st.Level
}

func (s *Service) currentAvatarKey(userID uint) string {
	user, _, _, _, err := s.repo.FindByID(userID)
	if err != nil || user.AvatarURL == nil {
		return ""
	}
	if strings.HasPrefix(*user.AvatarURL, "avatars/") {
		return *user.AvatarURL
	}
	return ""
}

func toProfileResponse(u *User, reviewCount, followerCount, followingCount, level int) *UserProfileResponse {
	return &UserProfileResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		DisplayName:    u.DisplayName,
		Bio:            u.Bio,
		AvatarURL:      u.AvatarURL,
		Gender:         u.Gender,
		FavoriteGenres: u.FavoriteGenres,
		DateOfBirth:    u.DateOfBirth,
		ReviewCount:    reviewCount,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		IsPrivate:      u.IsPrivate,
		Level:          level,
		Role:           string(u.Role.RoleName),
		CreatedAt:      u.CreatedAt,
	}
}
