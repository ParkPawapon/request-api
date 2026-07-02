package petitiontypes

import domain "github.com/ParkPawapon/request-api/internal/domain/petitiontype"

type PetitionTypeResponse struct {
	PetitionTypeID   int    `json:"petitionTypeID"`
	PetitionTypeName string `json:"petitionTypeName"`
}

func toResponse(items []domain.PetitionType) []PetitionTypeResponse {
	out := make([]PetitionTypeResponse, 0, len(items))
	for _, item := range items {
		out = append(out, PetitionTypeResponse{
			PetitionTypeID:   item.ID,
			PetitionTypeName: item.Name,
		})
	}
	return out
}
