package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type selfAssessmentRepository struct {
	db *gorm.DB
}

func NewSelfAssessmentRepository(db *gorm.DB) usecase.SelfAssessmentRepository {
	return &selfAssessmentRepository{db}
}

func (s *selfAssessmentRepository) CreateAssessmentSession(ctx context.Context, skills []entity.SelfAssessmentSkill) error {
	return s.db.WithContext(ctx).Create(&skills).Error
}
