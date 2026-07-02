package repository

import (
	"context"
	"fmt"

	"github.com/ParkPawapon/request-api/internal/domain/petitiontype"
	"gorm.io/gorm"
)

type PetitionTypeGormRepository struct {
	db *gorm.DB
}

func NewPetitionTypeGormRepository(db *gorm.DB) *PetitionTypeGormRepository {
	return &PetitionTypeGormRepository{db: db}
}

func (r *PetitionTypeGormRepository) ListActive(ctx context.Context) ([]petitiontype.PetitionType, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("petition type repository database is nil")
	}

	var rows []petitionTypeRecord
	if err := r.db.WithContext(ctx).
		Model(&petitionTypeRecord{}).
		Select(`"petitionTypeID", "petitionTypeName"`).
		Where(`"petitionTypeName" IS NOT NULL`).
		Order(`"petitionTypeName" ASC`).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]petitiontype.PetitionType, 0, len(rows))
	for _, row := range rows {
		items = append(items, petitiontype.PetitionType{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	return items, nil
}

type petitionTypeRecord struct {
	ID   int    `gorm:"column:petitionTypeID"`
	Name string `gorm:"column:petitionTypeName"`
}

func (petitionTypeRecord) TableName() string {
	return "petitionType"
}
