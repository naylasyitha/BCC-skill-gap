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

func (s *selfAssessmentRepository) CreateAssessmentSession(ctx context.Context, session *entity.UserCareerSession, skills []entity.SelfAssessmentSkill) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Create(session).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := range skills {
		skills[i].UserCareerSessionID = session.ID
	}

	err = tx.Create(&skills).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}
