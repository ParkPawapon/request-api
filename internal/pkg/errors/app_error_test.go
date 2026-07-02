package errors

import (
	stderrors "errors"
	"net/http"
	"testing"
)

func TestAppErrorStatusAndCode(t *testing.T) {
	err := Unauthorized("Authentication is required.", stderrors.New("missing session"))

	if err.Status != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, err.Status)
	}
	if err.Code != CodeUnauthorized {
		t.Fatalf("expected code %s, got %s", CodeUnauthorized, err.Code)
	}
	if err.Message != "Authentication is required." {
		t.Fatalf("unexpected message %q", err.Message)
	}
}

func TestAppErrorUnwrap(t *testing.T) {
	internal := stderrors.New("database connection refused")
	err := Internal("Internal Server Error", internal)

	if !stderrors.Is(err, internal) {
		t.Fatal("expected wrapped internal error")
	}
}
