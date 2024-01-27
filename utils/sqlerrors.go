package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

type SQLFieldError map[string][]string

func (e SQLFieldError) Error() string {
	var errStrings []string
	for field, errs := range e {
		for _, err := range errs {
			errStrings = append(errStrings, fmt.Sprintf("%s: %s", field, err))
		}
	}
	return strings.Join(errStrings, "; ")
}

func NewSQLFieldErrors(dupErr *pgconn.PgError) error {
	errors := make(SQLFieldError)
	if dupErr.Code == "23505" {
		// Extract field name from the error message
		// The constraint name should follow the format: 'unq-<table-name>-<fieldname>'
		parts := strings.Split(dupErr.ConstraintName, "-")

		if len(parts) != 3 {
			return fmt.Errorf("unexpected constraint name format: %s", dupErr.ConstraintName)
		}

		field := parts[2]
		errors[field] = append(errors[field], fmt.Sprintf("%s is already in use", field))
	}
	return errors
}
