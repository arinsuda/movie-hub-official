package movie_module

import (
	"math/rand/v2"
	"strconv"

	tmdb "github.com/arinsuda/movie-hub/internal/tmdb_module"
	"github.com/gofiber/fiber/v3"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetPopular(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := tmdb.GetPopular(page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลหนังไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetNowPlaying(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := tmdb.GetNowPlaying(page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลหนังไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetTopRated(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := tmdb.GetTopRated(page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลหนังไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetUpcoming(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := tmdb.GetUpcoming(page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลหนังไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) Search(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "กรุณาระบุคำค้นหา"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := tmdb.SearchMovies(query, page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ค้นหาไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetGenres(c fiber.Ctx) error {
	genres, err := tmdb.GetGenres()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูล genre ไม่สำเร็จ"})
	}
	return c.JSON(fiber.Map{"genres": genres})
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}

	movie, err := tmdb.GetMovieByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบหนัง"})
	}

	credits, _ := tmdb.GetCredits(id)
	videos, _ := tmdb.GetVideos(id)

	return c.JSON(fiber.Map{
		"movie":   movie,
		"credits": credits,
		"videos":  videos,
	})
}

func (h *Handler) GetSimilar(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := tmdb.GetSimilar(id, page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetPopularSeries(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	result, err := tmdb.GetPopularSeries(page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลซีรีส์ไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetNowAiringSeries(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	result, err := tmdb.GetNowAiringSeries(page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลซีรีส์ไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetTopRatedSeries(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	result, err := tmdb.GetTopRatedSeries(page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลซีรีส์ไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) SearchSeries(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "กรุณาระบุคำค้นหา"})
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	result, err := tmdb.SearchSeries(query, page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ค้นหาไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetSeriesGenres(c fiber.Ctx) error {
	genres, err := tmdb.GetSeriesGenres()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูล genre ไม่สำเร็จ"})
	}
	return c.JSON(fiber.Map{"genres": genres})
}

func (h *Handler) GetSeriesByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}

	series, err := tmdb.GetSeriesByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบซีรีส์"})
	}

	credits, _ := tmdb.GetSeriesCredits(id)
	videos, _ := tmdb.GetSeriesVideos(id)

	return c.JSON(fiber.Map{
		"series":  series,
		"credits": credits,
		"videos":  videos,
	})
}

func (h *Handler) GetSimilarSeries(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID ไม่ถูกต้อง"})
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	result, err := tmdb.GetSimilarSeries(id, page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ดึงข้อมูลไม่สำเร็จ"})
	}
	return c.JSON(result)
}

func (h *Handler) GetRecommended(c fiber.Ctx) error {
	genreIDs := c.Query("with_genres")
	if genreIDs == "" {
		return h.GetPopular(c)
	}

	// ดึงหลาย page แล้วรวมกัน เพื่อให้ได้หนังมากพอ
	var allResults []tmdb.Movie

	for page := 1; page <= 3; page++ {
		result, err := tmdb.DiscoverMovies(genreIDs, page)
		if err != nil {
			break
		}
		allResults = append(allResults, result.Results...)
		if page >= result.TotalPages {
			break
		}
	}

	rand.Shuffle(len(allResults), func(i, j int) {
		allResults[i], allResults[j] = allResults[j], allResults[i]
	})

	if len(allResults) > 10 {
		allResults = allResults[:10]
	}

	return c.JSON(fiber.Map{
		"results":       allResults,
		"total_results": len(allResults),
	})
}
