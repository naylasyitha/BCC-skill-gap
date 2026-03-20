package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type careerSkillRepository struct {
	db *gorm.DB
}

func NewCareerSkillRepository(db *gorm.DB) usecase.CareerSkillRepository {
	return &careerSkillRepository{db}
}

// Delete implements [usecase.CareerSkillRepository].
func (c *careerSkillRepository) Delete(ctx context.Context, id string) error {
	return c.db.WithContext(ctx).Delete(&entity.CareerSkill{}, id).Error
}

// FindById implements [usecase.CareerSkillRepository].
func (c *careerSkillRepository) FindById(ctx context.Context, id string) (*entity.CareerSkill, error) {
	var careerSkill entity.CareerSkill
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&careerSkill).Error
	if err != nil {
		return nil, err
	}

	return &careerSkill, nil
}

// Save implements [usecase.CareerSkillRepository].
func (c *careerSkillRepository) Save(ctx context.Context, careerSkill *entity.CareerSkill) error {
	return c.db.WithContext(ctx).Create(careerSkill).Error
}

// Update implements [usecase.CareerSkillRepository].
func (c *careerSkillRepository) Update(ctx context.Context, careerSkill *entity.CareerSkill) error {
	return c.db.WithContext(ctx).Save(careerSkill).Error
}
