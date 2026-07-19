package shared

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const PgUniqueViolation = "23505"

func IsPgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == PgUniqueViolation
	}
	return false
}

func IsPgUniqueViolationOn(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == PgUniqueViolation && pgErr.ConstraintName == constraint
	}
	return false
}
