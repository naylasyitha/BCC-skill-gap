package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type skillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) usecase.SkillRepository {
	return &skillRepository{db}
}

// Delete implements [usecase.SkillRepository].
func (s *skillRepository) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Skill{}).Error
}

// FindAll implements [usecase.SkillRepository].
func (s *skillRepository) FindAll(ctx context.Context) ([]entity.Skill, error) {
	var skill []entity.Skill
	err := s.db.WithContext(ctx).Find(&skill).Error
	if err != nil {
		return nil, err
	}

	return skill, nil
}

// FindById implements [usecase.SkillRepository].
func (s *skillRepository) FindById(ctx context.Context, id string) (*entity.Skill, error) {
	var skill entity.Skill
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&skill).Error
	if err != nil {
		return nil, err
	}

	return &skill, nil
}

// Save implements [usecase.SkillRepository].
func (s *skillRepository) Save(ctx context.Context, skill *entity.Skill) error {
	return s.db.WithContext(ctx).Create(skill).Error
}

// Update implements [usecase.SkillRepository].
func (s *skillRepository) Update(ctx context.Context, skill *entity.Skill) error {
	return s.db.WithContext(ctx).Save(skill).Error
}
