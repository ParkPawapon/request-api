package repository

import (
	"context"

	"github.com/ParkPawapon/request-api/internal/domain/petitiontype"
)

type PetitionTypeRepository interface {
	ListActive(ctx context.Context) ([]petitiontype.PetitionType, error)
}
