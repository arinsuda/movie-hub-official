package auth_module

import (
	"errors"
	"time"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/user_module"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
	cfg *config.Config
}

func NewHandler(svc *Service, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}

func (h *Handler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := h.svc.Register(req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		User: toUserResponse(user),
	})
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	pair, user, err := h.svc.Login(req, c.Get("User-Agent"), c.IP())
	if err != nil {
		return h.handleServiceError(c, err)
	}

	h.setTokenCookies(c, pair)

	return c.JSON(AuthResponse{User: toUserResponse(user)})
}

func (h *Handler) Refresh(c fiber.Ctx) error {
	rawRefresh := c.Cookies("refresh_token")
	if rawRefresh == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing refresh token"})
	}

	pair, user, err := h.svc.Refresh(rawRefresh, c.Get("User-Agent"), c.IP())
	if err != nil {
		h.clearTokenCookies(c)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid session"})
	}

	h.setTokenCookies(c, pair)

	return c.JSON(AuthResponse{User: toUserResponse(user)})
}

func (h *Handler) Logout(c fiber.Ctx) error {
	rawRefresh := c.Cookies("refresh_token")
	if rawRefresh != "" {
		_ = h.svc.Logout(rawRefresh)
	}
	h.clearTokenCookies(c)
	return c.JSON(MessageResponse{Message: "logged out successfully"})
}

func (h *Handler) LogoutAll(c fiber.Ctx) error {
	claims := GetClaims(c)
	_ = h.svc.LogoutAll(claims.UserID)
	h.clearTokenCookies(c)
	return c.JSON(MessageResponse{Message: "logged out from all devices"})
}

func (h *Handler) VerifyEmail(c fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing token"})
	}

	if err := h.svc.VerifyEmail(token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	return c.JSON(MessageResponse{Message: "email verified successfully"})
}

func (h *Handler) ResendVerification(c fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.Bind().JSON(&body); err != nil || body.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email required"})
	}

	if err := h.svc.ResendVerification(body.Email); err != nil {
		if errors.Is(err, ErrAlreadyVerified) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already verified"})
		}
	}

	return c.JSON(MessageResponse{Message: "if the email exists, a verification link has been sent"})
}

func (h *Handler) setTokenCookies(c fiber.Ctx, pair *TokenPair) {
	sameSite := h.cfg.Cookie.SameSite

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    pair.AccessToken,
		HTTPOnly: true,
		Secure:   h.cfg.Cookie.Secure,
		SameSite: sameSite,
		Domain:   h.cfg.Cookie.Domain,
		MaxAge:   int(h.cfg.JWT.AccessTTL.Seconds()),
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		HTTPOnly: true,
		Secure:   h.cfg.Cookie.Secure,
		SameSite: sameSite,
		Domain:   h.cfg.Cookie.Domain,
		MaxAge:   int(h.cfg.JWT.RefreshTTL.Seconds()),
		Path:     "/auth/refresh",
	})
}

func (h *Handler) clearTokenCookies(c fiber.Ctx) {
	for _, name := range []string{"access_token", "refresh_token"} {
		c.Cookie(&fiber.Cookie{
			Name:    name,
			Value:   "",
			MaxAge:  -1,
			Expires: time.Unix(0, 0),
			Path:    "/",
		})
	}
}

func (h *Handler) handleServiceError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrPasswordMismatch):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, ErrEmailTaken), errors.Is(err, ErrUsernameTaken):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, ErrUserNotFound), errors.Is(err, ErrWrongPassword):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	case errors.Is(err, ErrEmailUnverified):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "please verify your email first"})
	case errors.Is(err, ErrInvalidToken):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

func toUserResponse(u *user_module.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		IsVerified:  u.VerifiedEmailAt != nil,
		Role:        string(u.Role.RoleName),
	}
}
