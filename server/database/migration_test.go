package database

import (
	"testing"
)

func TestHashString(t *testing.T) {
	h1 := hashString("test_migration")
	h2 := hashString("test_migration")
	if h1 != h2 {
		t.Errorf("expected hashes to be equal, got %d and %d", h1, h2)
	}
}

func TestErrActiveDuplicateReviews(t *testing.T) {
	dupErr := &ErrActiveDuplicateReviews{
		Groups: []DuplicateReviewGroup{
			{UserID: 1, MediaID: 42, MediaType: "movie", ReviewIDs: []uint{10, 11}},
		},
	}
	if dupErr.Error() != "active duplicate reviews detected" {
		t.Errorf("unexpected error message: %s", dupErr.Error())
	}
}
