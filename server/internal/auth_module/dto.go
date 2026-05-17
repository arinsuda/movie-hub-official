package auth_module

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email           string `json:"email"    validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,eqfield=Password"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password"   validate:"required"`
}

type UserResponse struct {
	ID             uint    `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	DisplayName    *string `json:"display_name"`
	AvatarURL      *string `json:"avatar_url"`
	IsVerified     bool    `json:"is_verified"`
	Role           string  `json:"role"`
	FavoriteGenres *string `json:"favorite_genres"`
}

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
