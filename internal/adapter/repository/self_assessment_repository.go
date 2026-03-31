package repository

import (
	"context"
	"errors"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"github.com/google/uuid"
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

func (s *selfAssessmentRepository) UpdateStatus(ctx context.Context, careerSessionID string, status entity.StatusEnum) error {
	careerUUID, err := uuid.Parse(careerSessionID)
	if err != nil {
		return errors.New("Format ID tidak valid")
	}
	return s.db.WithContext(ctx).Model(&entity.UserCareerSession{}).Where("id = ?", careerUUID).Update("status", status).Error
}
