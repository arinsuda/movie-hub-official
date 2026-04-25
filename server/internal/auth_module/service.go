package auth_module

import (
	"fmt"
	"time"

	"github.com/arinsuda/movie-hub/internal/config"
	"github.com/arinsuda/movie-hub/internal/mailer"
	"github.com/arinsuda/movie-hub/internal/user_module"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo   *repository
	jwt    *jwtManager
	mailer *mailer.Mailer
	cfg    *config.Config
}

func NewService(db *gorm.DB, cfg *config.Config, m *mailer.Mailer) *Service {
	return &Service{
		repo:   newRepository(db),
		jwt:    newJWTManager(cfg.JWT),
		mailer: m,
		cfg:    cfg,
	}
}

// ── Register ──────────────────────────────────────────────────────

func (s *Service) Register(req RegisterRequest) (*user_module.User, error) {
	// ตรวจ duplicate
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, ErrEmailTaken
	}
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return nil, ErrUsernameTaken
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}

	user := &user_module.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
		RoleID:   2, // user role
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("register: create user: %w", err)
	}

	// ส่ง verification email (async ไม่บล็อก response)
	go func() {
		_ = s.sendVerificationEmail(user)
	}()

	return user, nil
}

// ── Login ─────────────────────────────────────────────────────────

func (s *Service) Login(req LoginRequest, userAgent, ip string) (*TokenPair, *user_module.User, error) {
	// หา user จาก email หรือ username
	user, err := s.repo.FindByEmail(req.Identifier)
	if err != nil {
		user, err = s.repo.FindByUsername(req.Identifier)
		if err != nil {
			return nil, nil, ErrUserNotFound
		}
	}

	if !user.IsActive {
		return nil, nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, nil, ErrWrongPassword
	}

	if user.VerifiedEmailAt == nil {
		return nil, nil, ErrEmailUnverified
	}

	pair, err := s.issueAndStoreTokens(user, userAgent, ip)
	if err != nil {
		return nil, nil, err
	}

	return pair, user, nil
}

// ── Refresh ───────────────────────────────────────────────────────

func (s *Service) Refresh(rawRefreshToken, userAgent, ip string) (*TokenPair, *user_module.User, error) {
	claims, err := s.jwt.ParseRefresh(rawRefreshToken)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	hashed := HashToken(rawRefreshToken)

	stored, err := s.repo.FindRefreshToken(hashed)
	if err != nil || stored.IsExpired() {
		return nil, nil, ErrInvalidToken
	}

	// Rotate: ลบ token เก่า
	if err := s.repo.DeleteRefreshToken(hashed); err != nil {
		return nil, nil, fmt.Errorf("refresh: delete old token: %w", err)
	}

	user, err := s.repo.FindByID(claims.UserID)
	if err != nil || !user.IsActive {
		return nil, nil, ErrUserNotFound
	}

	pair, err := s.issueAndStoreTokens(user, userAgent, ip)
	if err != nil {
		return nil, nil, err
	}

	return pair, user, nil
}

// ── Logout ────────────────────────────────────────────────────────

func (s *Service) Logout(rawRefreshToken string) error {
	hashed := HashToken(rawRefreshToken)
	return s.repo.DeleteRefreshToken(hashed)
}

func (s *Service) LogoutAll(userID uint) error {
	return s.repo.DeleteAllRefreshTokens(userID)
}

// ── Email Verification ────────────────────────────────────────────

func (s *Service) VerifyEmail(rawToken string) error {
	hashed := HashToken(rawToken)

	ev, err := s.repo.FindEmailVerification(hashed)
	if err != nil {
		return ErrInvalidToken
	}
	if ev.IsExpired() {
		_ = s.repo.DeleteEmailVerification(ev.ID)
		return ErrInvalidToken
	}

	if err := s.repo.MarkEmailVerified(ev.UserID); err != nil {
		return fmt.Errorf("verify email: update user: %w", err)
	}

	return s.repo.DeleteEmailVerification(ev.ID)
}

func (s *Service) ResendVerification(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// ไม่บอกว่า email ไม่มีในระบบ (security)
		return nil
	}
	if user.VerifiedEmailAt != nil {
		return ErrAlreadyVerified
	}

	return s.sendVerificationEmail(user)
}

// ── internal helpers ──────────────────────────────────────────────

func (s *Service) issueAndStoreTokens(user *user_module.User, userAgent, ip string) (*TokenPair, error) {
	pair, err := s.jwt.Issue(user.ID, string(user.Role.RoleName))
	if err != nil {
		return nil, fmt.Errorf("issue tokens: %w", err)
	}

	rt := &user_module.RefreshToken{
		UserID:      user.ID,
		HashedToken: HashToken(pair.RefreshToken),
		ExpiresAt:   time.Now().Add(s.cfg.JWT.RefreshTTL),
		UserAgent:   userAgent,
		IPAddress:   ip,
	}
	if err := s.repo.CreateRefreshToken(rt); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return pair, nil
}

func (s *Service) sendVerificationEmail(user *user_module.User) error {
	rawToken, err := GenerateSecureToken()
	if err != nil {
		return fmt.Errorf("send verify: generate token: %w", err)
	}

	ev := &user_module.EmailVerification{
		UserID:    user.ID,
		Token:     HashToken(rawToken),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.repo.CreateEmailVerification(ev); err != nil {
		return fmt.Errorf("send verify: save token: %w", err)
	}

	verifyURL := fmt.Sprintf("%s/auth/verify-email?token=%s",
		s.cfg.AppBaseURL, rawToken)

	return s.mailer.SendVerificationEmail(user.Email, user.Username, verifyURL)
}
