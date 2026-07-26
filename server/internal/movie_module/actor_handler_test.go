package movie_module

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestActorRoutesAndHandlers(t *testing.T) {
	app := fiber.New()
	svc := NewMovieService(nil)
	h := NewHandler(svc)

	actors := app.Group("/actors")
	actors.Get("/search", h.SearchActor)
	actors.Get("/:id", h.GetActorDetails)
	actors.Get("/:id/movies", h.GetMoviesByActor)
	actors.Get("/:id/tv", h.GetSeriesByActor)

	t.Run("GetActorDetails invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/actors/abc", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("GetMoviesByActor invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/actors/invalid/movies", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("GetSeriesByActor invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/actors/invalid/tv", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("SearchActor empty query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/actors/search", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}
