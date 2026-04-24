package tmdbmodule

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const defaultLanguage = "th-TH"

func ImageURL(path string) string {
	if path == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w500" + path
}

type client struct {
	http     *http.Client
	baseURL  string
	apiKey   string
	language string
}

var c *client

func Init() {
	c = &client{
		http:     &http.Client{Timeout: 10 * time.Second},
		baseURL:  os.Getenv("THE_MOVIE_BASE_API"),
		apiKey:   os.Getenv("THE_MOVIE_API_KEY"),
		language: defaultLanguage,
	}
}

func get(path string, params url.Values, target any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_key", c.apiKey)
	if params.Get("language") == "" {
		params.Set("language", c.language)
	}
	fullURL := fmt.Sprintf("%s%s?%s", c.baseURL, path, params.Encode())

	resp, err := c.http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("tmdb: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("tmdb: status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func getWithFallback(path string, params url.Values, target any, isEmpty func() bool) error {
	if err := get(path, params, target); err != nil {
		return err
	}

	if isEmpty() {
		if params == nil {
			params = url.Values{}
		}
		params.Set("language", "en-US")
		return get(path, params, target)
	}

	return nil
}
