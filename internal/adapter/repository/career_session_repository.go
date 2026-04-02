package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

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
	err := c.db.WithContext(ctx).
		Preload("User").
		Preload("Career").
		First(&careerSession, "id = ?", careerSessionId).Error

	if err != nil {
		return nil, err
	}
	return &careerSession, nil
}
func (c *careerSessionRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	var count int64
	err := c.db.WithContext(ctx).
		Model(&entity.UserCareerSession{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

func (c *careerSessionRepository) GetAllCareerSession(ctx context.Context, userID string) ([]entity.UserCareerSession, error) {
	var careerSessions []entity.UserCareerSession
	err := c.db.WithContext(ctx).Preload("Career").Where("user_id = ?", userID).Find(&careerSessions).Error
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

func (c *careerSessionRepository) GetRequiredLevel(ctx context.Context, careerSessionID string) ([]entity.CareerSkill, error) {
	var userCareerSession entity.UserCareerSession
	err := c.db.WithContext(ctx).Select("career_id").Where("id = ?", careerSessionID).First(&userCareerSession).Error
	if err != nil {
		return nil, err
	}

	var requiredLevels []entity.CareerSkill
	err = c.db.WithContext(ctx).Where("career_id", userCareerSession.CareerID).Find(&requiredLevels).Error
	if err != nil {
		return nil, err
	}

	return requiredLevels, nil
}
