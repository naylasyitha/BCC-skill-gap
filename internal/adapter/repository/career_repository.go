package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type careerRepository struct {
	db *gorm.DB
}

func NewCareerRepository(db *gorm.DB) usecase.CareerRepository {
	return &careerRepository{db}
}

// Delete implements [usecase.CareerRepository].
func (c *careerRepository) Delete(ctx context.Context, id string) error {
	return c.db.WithContext(ctx).Delete(&entity.Career{}, id).Error
}

// FindAll implements [usecase.CareerRepository].
func (c *careerRepository) FindAll(ctx context.Context) ([]entity.Career, error) {
	var career []entity.Career
	err := c.db.WithContext(ctx).Find(&career).Error
	if err != nil {
		return nil, err
	}
	return career, nil
}

// FindById implements [usecase.CareerRepository].
func (c *careerRepository) FindById(ctx context.Context, id string) (*entity.Career, error) {
	var career entity.Career
	err := c.db.WithContext(ctx).Preload("CareerSkills.Skill").Where("id = ?", id).First(&career).Error
	if err != nil {
		return nil, err
	}
	return &career, nil
}

// Save implements [usecase.CareerRepository].
func (c *careerRepository) Save(ctx context.Context, career *entity.Career) error {
	return c.db.WithContext(ctx).Create(career).Error
}

// Update implements [usecase.CareerRepository].
func (c *careerRepository) Update(ctx context.Context, career *entity.Career) error {
	return c.db.WithContext(ctx).Save(career).Error
}

func (c *careerRepository) CreateCareerSkill(ctx context.Context, career *entity.Career, careerSkill []entity.CareerSkill) error {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(career).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range careerSkill {
		careerSkill[i].CareerID = career.ID
	}

	if err := tx.Create(&careerSkill).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
