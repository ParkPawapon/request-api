package security

import "errors"

var ErrPasswordAuthNotMigrated = errors.New("password authentication is not migrated")

func HashPassword(_ string) (string, error) {
	return "", ErrPasswordAuthNotMigrated
}

func ComparePassword(_, _ string) (bool, error) {
	return false, ErrPasswordAuthNotMigrated
}
