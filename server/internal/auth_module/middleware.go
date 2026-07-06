package auth_module

import (
	"github.com/arinsuda/movie-hub/config"
	mw "github.com/arinsuda/movie-hub/middleware"
)

type Middleware = mw.AuthMiddleware

func NewMiddleware(cfg *config.Config) *mw.AuthMiddleware {
	return mw.NewAuthMiddleware(cfg)
}

var GetClaims = mw.GetClaims
