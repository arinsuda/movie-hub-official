package shared

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsPgUniqueViolationOn(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:           PgUniqueViolation,
		ConstraintName: "uq_active_user_media_review",
	}

	if !IsPgUniqueViolationOn(pgErr, "uq_active_user_media_review") {
		t.Error("expected true for matching pgconn unique violation error")
	}

	if IsPgUniqueViolationOn(pgErr, "other_constraint") {
		t.Error("expected false for mismatched constraint name")
	}

	standardErr := errors.New("some standard error")
	if IsPgUniqueViolationOn(standardErr, "uq_active_user_media_review") {
		t.Error("expected false for standard non-pgconn error")
	}
}

func TestIsPgUniqueViolation(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code: PgUniqueViolation,
	}

	if !IsPgUniqueViolation(pgErr) {
		t.Error("expected true for pgconn unique violation error")
	}

	standardErr := errors.New("some standard error")
	if IsPgUniqueViolation(standardErr) {
		t.Error("expected false for standard non-pgconn error")
	}
}

func TestIsValidRating(t *testing.T) {
	validRatings := []float32{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0}
	for _, rating := range validRatings {
		if !IsValidRating(rating) {
			t.Errorf("expected rating %f to be valid", rating)
		}
	}

	invalidRatings := []float32{0.0, 0.4, 0.6, 2.3, 4.9, 5.1, 10.0}
	for _, rating := range invalidRatings {
		if IsValidRating(rating) {
			t.Errorf("expected rating %f to be invalid", rating)
		}
	}
}
