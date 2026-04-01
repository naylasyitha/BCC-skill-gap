package repository

import (
	"context"
	"errors"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type careerSessionRepository struct {
	db *gorm.DB
}

func NewCareerSessionRepository(db *gorm.DB) usecase.CareerSessionRepository {
	return &careerSessionRepository{db}
}

func (c *careerSessionRepository) Create(ctx context.Context, session *entity.UserCareerSession) error {
	return c.db.WithContext(ctx).Create(session).Error
}

func (c *careerSessionRepository) FindById(ctx context.Context, careerSessionId string) (*entity.UserCareerSession, error) {
	var careerSession entity.UserCareerSession
	careerUUID, err := uuid.Parse(careerSessionId)
	if err != nil {
		return nil, errors.New("Format ID tidak valid")
	}

	err = c.db.WithContext(ctx).Preload("User").Preload("Career").Where("id = ?", careerUUID).First(&careerSession).Error
	if err != nil {
		return nil, err
	}

	return &careerSession, nil
}
func (c *careerSessionRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, errors.New("Format ID tidak valid")
	}

	var count int64
	err = c.db.WithContext(ctx).Model(&entity.UserCareerSession{}).Where("user_id = ?", userUUID).Count(&count).Error
	return int(count), err
}

func (c *careerSessionRepository) GetAllCareerSession(ctx context.Context, userID string) ([]entity.UserCareerSession, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("Format ID tidak valid")
	}

	var careerSessions []entity.UserCareerSession
	err = c.db.WithContext(ctx).Preload("Career").Where("user_id = ?", userUUID).Find(&careerSessions).Error
	if err != nil {
		return nil, err
	}

	return careerSessions, nil
}

func (c *careerSessionRepository) GetAnalyticsData(ctx context.Context, careerSessionID string) ([]entity.SelfAssessmentSkill, error) {
	var skills []entity.SelfAssessmentSkill
	err := c.db.WithContext(ctx).Preload("Skill").Where("user_career_session_id = ?", careerSessionID).Find(&skills).Error
	if err != nil {
		return nil, err
	}
	return skills, nil
}
