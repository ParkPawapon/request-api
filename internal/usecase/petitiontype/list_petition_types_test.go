package petitiontype

import (
	"context"
	"errors"
	"testing"

	domain "github.com/ParkPawapon/request-api/internal/domain/petitiontype"
	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
)

func TestListPetitionTypesUseCaseReturnsRepositoryItems(t *testing.T) {
	useCase := NewListPetitionTypesUseCase(fakeRepository{
		items: []domain.PetitionType{
			{ID: 2, Name: "ขอลาออก"},
			{ID: 1, Name: "คำร้องทั่วไป"},
		},
	})

	items, err := useCase.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].ID != 2 || items[0].Name != "ขอลาออก" {
		t.Fatalf("unexpected first item: %+v", items[0])
	}
}

func TestListPetitionTypesUseCaseNormalizesRepositoryError(t *testing.T) {
	useCase := NewListPetitionTypesUseCase(fakeRepository{
		err: errors.New("postgres internal detail"),
	})

	_, err := useCase.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != apperrors.CodeInternal || appErr.Message != "Internal Server Error" {
		t.Fatalf("unexpected app error: %+v", appErr)
	}
}

func TestListPetitionTypesUseCaseFailsClosedWithoutRepository(t *testing.T) {
	useCase := NewListPetitionTypesUseCase(nil)

	_, err := useCase.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != apperrors.CodeServiceNotReady {
		t.Fatalf("expected service not ready, got %s", appErr.Code)
	}
}

type fakeRepository struct {
	items []domain.PetitionType
	err   error
}

func (r fakeRepository) ListActive(context.Context) ([]domain.PetitionType, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.items, nil
}
