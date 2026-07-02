package security

import "errors"

var ErrTokenAuthNotMigrated = errors.New("token authentication is not migrated")

type TokenClaims struct {
	Subject string
}

func IssueToken(_ TokenClaims) (string, error) {
	return "", ErrTokenAuthNotMigrated
}

func ParseToken(_ string) (TokenClaims, error) {
	return TokenClaims{}, ErrTokenAuthNotMigrated
}
