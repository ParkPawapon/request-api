package petitiontype

import (
	"context"

	domain "github.com/ParkPawapon/request-api/internal/domain/petitiontype"
	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/repository"
)

type ListPetitionTypesUseCase struct {
	repo repository.PetitionTypeRepository
}

func NewListPetitionTypesUseCase(repo repository.PetitionTypeRepository) *ListPetitionTypesUseCase {
	return &ListPetitionTypesUseCase{repo: repo}
}

func (u *ListPetitionTypesUseCase) Execute(ctx context.Context) ([]domain.PetitionType, error) {
	if u == nil || u.repo == nil {
		return nil, apperrors.ServiceNotReady("Petition type service is not ready.", nil)
	}

	items, err := u.repo.ListActive(ctx)
	if err != nil {
		return nil, apperrors.Internal("Internal Server Error", err)
	}

	return items, nil
}
