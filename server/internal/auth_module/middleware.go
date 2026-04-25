package auth_module

import (
	"github.com/arinsuda/movie-hub/config"
	mw "github.com/arinsuda/movie-hub/middleware"
)

// re-export เพื่อ backward compat (optional)
type Middleware = mw.AuthMiddleware

func NewMiddleware(cfg *config.Config) *mw.AuthMiddleware {
	return mw.NewAuthMiddleware(cfg)
}

// GetClaims ก็ re-export ได้
var GetClaims = mw.GetClaims
