package user_module

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	mw "github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProfile(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	profile, err := h.svc.GetProfile(targetID, claims.UserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"user": profile})
}

func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	req, err := parseUpdateProfileForm(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var avatarFile *multipartFile
	fh, formErr := c.FormFile("avatar")
	if formErr == nil && fh != nil {
		if err := validateImageFile(fh); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		f, err := fh.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot read avatar file"})
		}
		defer f.Close()
		avatarFile = &multipartFile{file: f, header: fh}
	}

	var profile *UserProfileResponse
	if avatarFile != nil {
		profile, err = h.svc.UpdateProfile(targetID, claims.UserID, req, avatarFile.file, avatarFile.header)
	} else {
		profile, err = h.svc.UpdateProfile(targetID, claims.UserID, req, nil, nil)
	}
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"user": profile})
}

func (h *Handler) UpdateEmail(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	claims := mw.GetClaims(c)

	var req UpdateEmailRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.NewEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new_email is required"})
	}

	profile, err := h.svc.UpdateEmail(targetID, claims.UserID, req.NewEmail)
	if err != nil {
		return handleEmailChangeError(c, err)
	}
	return c.JSON(fiber.Map{"user": profile})
}

func (h *Handler) DeleteUser(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	if err := h.svc.DeleteUser(targetID, claims.UserID, claims.Role); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *Handler) UpdateFavoriteGenres(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	var req UpdateFavoriteGenresRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	profile, err := h.svc.UpdateFavoriteGenres(targetID, claims.UserID, req.FavoriteGenres)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"user": profile})
}

func (h *Handler) RequestEmailChange(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	if err := h.svc.RequestEmailChange(targetID, claims.UserID); err != nil {
		return handleEmailChangeError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "OTP has been sent to your current email address",
	})
}

func (h *Handler) VerifyEmailChange(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	var req VerifyEmailChangeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "otp is required"})
	}

	profile, err := h.svc.VerifyEmailChange(targetID, claims.UserID, req.OTP)
	if err != nil {
		return handleEmailChangeError(c, err)
	}

	return c.JSON(fiber.Map{"user": profile})
}

// ── Password ──────────────────────────────────────────────────────────

// ChangePassword: Case 1 — จำรหัสผ่านเดิมได้
// PATCH /users/:userId/password
func (h *Handler) ChangePassword(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	var req ChangePasswordRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.OldPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "old_password is required"})
	}
	if req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new_password is required"})
	}
	if req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "confirm_password is required"})
	}

	if err := h.svc.ChangePassword(targetID, claims.UserID, req); err != nil {
		return handlePasswordError(c, err)
	}

	return c.JSON(fiber.Map{"message": "password changed successfully"})
}

// ForgotPassword: Case 2A — ขอ reset link (public endpoint ไม่ต้อง auth)
// POST /auth/forgot-password
func (h *Handler) ForgotPassword(c fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	// เรียก service แล้ว return 200 เสมอ (ป้องกัน user enumeration)
	_ = h.svc.ForgotPassword(req.Email)

	return c.JSON(fiber.Map{
		"message": "if an account with that email exists, a password reset link has been sent",
	})
}

// ResetPassword: Case 2B — ตั้งรหัสผ่านใหม่ด้วย token (public endpoint)
// POST /auth/reset-password
func (h *Handler) ResetPassword(c fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token is required"})
	}
	if req.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}
	if req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new_password is required"})
	}
	if req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "confirm_password is required"})
	}

	if err := h.svc.ResetPassword(req.UserID, req.Token, req); err != nil {
		return handlePasswordError(c, err)
	}

	return c.JSON(fiber.Map{"message": "password has been reset successfully"})
}

// ── error handlers ────────────────────────────────────────────────────

func handlePasswordError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	case errors.Is(err, ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	case errors.Is(err, ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect current password"})
	case errors.Is(err, ErrPasswordResetTokenNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid or expired reset link"})
	case errors.Is(err, ErrPasswordResetTokenExpired):
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "reset link has expired, please request a new one"})
	case errors.Is(err, ErrPasswordResetTokenInvalid):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid reset token"})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
}

func handleEmailChangeError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	case errors.Is(err, ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	case errors.Is(err, ErrOTPNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no pending email change request"})
	case errors.Is(err, ErrOTPExpired):
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "otp has expired, please request a new one"})
	case errors.Is(err, ErrOTPInvalid):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid otp"})
	case errors.Is(err, ErrOTPMaxAttempts):
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many failed attempts, please request a new otp"})
	case errors.Is(err, ErrEmailAlreadyInUse):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already in use"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
}

// ── form helpers ──────────────────────────────────────────────────────

type multipartFile struct {
	file   multipart.File
	header *multipart.FileHeader
}

func parseUpdateProfileForm(c fiber.Ctx) (UpdateProfileRequest, error) {
	allowed := map[string]bool{
		"display_name":    true,
		"bio":             true,
		"gender":          true,
		"gender_other":    true,
		"favorite_genres": true,
		"date_of_birth":   true,
		"is_private":      true,
	}

	form, err := c.MultipartForm()
	if err != nil {
		return UpdateProfileRequest{}, errors.New("invalid multipart form")
	}

	var unknownFields []string
	for key := range form.Value {
		if !allowed[key] {
			unknownFields = append(unknownFields, key)
		}
	}
	if len(unknownFields) > 0 {
		return UpdateProfileRequest{}, fmt.Errorf("unknown fields: %s", strings.Join(unknownFields, ", "))
	}

	var req UpdateProfileRequest
	if v := c.FormValue("display_name"); v != "" {
		req.DisplayName = &v
	}
	if v := c.FormValue("bio"); v != "" {
		req.Bio = &v
	}
	if v := c.FormValue("gender"); v != "" {
		req.Gender = GenderType(v)
	}
	if v := c.FormValue("gender_other"); v != "" {
		req.GenderOther = &v
	}
	if v := c.FormValue("favorite_genres"); v != "" {
		req.FavoriteGenres = &v
	}
	if v := c.FormValue("date_of_birth"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return req, errors.New("date_of_birth must be YYYY-MM-DD")
		}
		req.DateOfBirth = &t
	}
	if v := c.FormValue("is_private"); v != "" {
		b := v == "true"
		req.IsPrivate = &b
	}

	return req, nil
}

func validateImageFile(fh *multipart.FileHeader) error {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}
	if !allowed[fh.Header.Get("Content-Type")] {
		return errors.New("avatar must be jpeg, png, webp, or gif")
	}
	const maxSize = 5 << 20 // 5 MB
	if fh.Size > maxSize {
		return errors.New("avatar must be smaller than 5 MB")
	}
	return nil
}

func parseUserID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("userId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func handleError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	case errors.Is(err, ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
