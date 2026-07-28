package shared

import "github.com/gofiber/fiber/v3"

type ErrorCode string

const (
	ErrorCodeDuplicateReview   ErrorCode = "DUPLICATE_REVIEW"
	ErrorCodeInvalidRating     ErrorCode = "INVALID_RATING"
	ErrorCodeInvalidVisibility ErrorCode = "INVALID_VISIBILITY"
	ErrorCodeWatchLogNotFound  ErrorCode = "WATCH_LOG_NOT_FOUND"
	ErrorCodeInvalidDate       ErrorCode = "INVALID_DATE"
	ErrorCodeFutureDate        ErrorCode = "FUTURE_DATE"
)

type APIError struct {
	Code  ErrorCode `json:"code"`
	Error string    `json:"error"`
}

func WriteAPIError(c fiber.Ctx, status int, code ErrorCode, message string) error {
	return c.Status(status).JSON(APIError{Code: code, Error: message})
}
