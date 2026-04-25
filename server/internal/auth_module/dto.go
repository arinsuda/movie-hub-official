package auth_module

// ── Request DTOs ────────────────────────────────────────────────

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	// รับ username หรือ email ก็ได้
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password"   validate:"required"`
}

// ── Response DTOs ───────────────────────────────────────────────

type UserResponse struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
	IsVerified  bool    `json:"is_verified"`
	Role        string  `json:"role"`
}

type AuthResponse struct {
	User UserResponse `json:"user"`
	// tokens อยู่ใน cookie เท่านั้น ไม่ส่งใน body
}

type MessageResponse struct {
	Message string `json:"message"`
}
