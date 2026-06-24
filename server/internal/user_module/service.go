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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo          *repository
	minio         *storage.MinIOClient
	statsSvc      StatsProvider
	mailer        Mailer
	emailVerifier EmailVerificationSender
}

func NewService(db *gorm.DB, mc *storage.MinIOClient, statsSvc StatsProvider, mailer Mailer, emailVerifier EmailVerificationSender) *Service {
	return &Service{
		repo:          newRepository(db),
		minio:         mc,
		statsSvc:      statsSvc,
		mailer:        mailer,
		emailVerifier: emailVerifier,
	}
}

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

func (s *Service) UpdateEmail(targetUserID, requesterID uint, newEmail string) (*UserProfileResponse, error) {
	if targetUserID != requesterID {
		return nil, ErrForbidden
	}

	newEmail = strings.TrimSpace(strings.ToLower(newEmail))
	if !isValidEmail(newEmail) {
		return nil, errors.New("invalid email format")
	}

	user, _, _, _, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(user.Email, newEmail) {
		return nil, errors.New("new email must be different from current email")
	}
	taken, err := s.repo.IsEmailTaken(newEmail, targetUserID)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrEmailAlreadyInUse
	}

	if err := s.repo.UpdateProfile(targetUserID, map[string]any{
		"email":             newEmail,
		"verified_email_at": nil,
	}); err != nil {
		return nil, err
	}

	defer func() {
		_ = s.emailVerifier.SendVerification(targetUserID, newEmail)
	}()

	return s.GetProfile(targetUserID, requesterID)
}

func (s *Service) DeleteUser(targetUserID, requesterID uint, requesterRole string) error {
	if targetUserID != requesterID && requesterRole != "admin" {
		return ErrForbidden
	}
	return s.repo.DeleteUser(targetUserID)
}

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

func (s *Service) RequestEmailChange(targetUserID, requesterID uint) error {
	if targetUserID != requesterID {
		return ErrForbidden
	}

	user, _, _, _, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return err
	}

	otp, err := generateOTP()
	if err != nil {
		return err
	}

	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	changeReq := &EmailChangeRequest{
		UserID:    targetUserID,
		OTPHash:   string(otpHash),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.repo.UpsertEmailChangeRequest(changeReq); err != nil {
		return err
	}

	if err := s.mailer.SendOTP(user.Email, otp); err != nil {

		_ = s.repo.DeleteEmailChangeRequest(targetUserID)
		return errors.New("failed to send OTP email")
	}

	return nil
}

const maxOTPAttempts = 5

func (s *Service) VerifyEmailChange(
	targetUserID, requesterID uint,
	otp string,
) (*UserProfileResponse, error) {
	if targetUserID != requesterID {
		return nil, ErrForbidden
	}

	req, err := s.repo.FindEmailChangeRequest(targetUserID)
	if err != nil {
		return nil, err
	}

	if req.IsExpired() {
		_ = s.repo.DeleteEmailChangeRequest(targetUserID)
		return nil, ErrOTPExpired
	}

	if req.AttemptCount >= maxOTPAttempts {
		_ = s.repo.DeleteEmailChangeRequest(targetUserID)
		return nil, ErrOTPMaxAttempts
	}

	if err := bcrypt.CompareHashAndPassword([]byte(req.OTPHash), []byte(otp)); err != nil {
		_ = s.repo.IncrementOTPAttempt(req.ID)
		return nil, ErrOTPInvalid
	}

	_ = s.repo.DeleteEmailChangeRequest(targetUserID)

	return s.GetProfile(targetUserID, requesterID)
}

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func isValidEmail(email string) bool {

	const pattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func (s *Service) fetchLevel(userID uint) int {
	if s.statsSvc == nil {
		return 1
	}
	return s.statsSvc.GetLevel(userID)
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
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		VerifiedEmailAt: u.VerifiedEmailAt,
		DisplayName:     u.DisplayName,
		Bio:             u.Bio,
		AvatarURL:       u.AvatarURL,
		Gender:          u.Gender,
		FavoriteGenres:  u.FavoriteGenres,
		DateOfBirth:     u.DateOfBirth,
		ReviewCount:     reviewCount,
		FollowerCount:   followerCount,
		FollowingCount:  followingCount,
		IsPrivate:       u.IsPrivate,
		Level:           level,
		Role:            string(u.Role.RoleName),
		CreatedAt:       u.CreatedAt,
	}
}
