package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) usecase.QuestionRepository {
	return &questionRepository{db}
}

func (q *questionRepository) Create(ctx context.Context, question *entity.Question) error {
	return q.db.WithContext(ctx).Create(question).Error
}

func (s *questionRepository) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Question{}).Error
}

func (s *questionRepository) FindAll(ctx context.Context) ([]entity.Question, error) {
	var question []entity.Question
	err := s.db.WithContext(ctx).Find(&question).Error
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (s *questionRepository) FindById(ctx context.Context, id string) (*entity.Question, error) {
	var question entity.Question
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (s *questionRepository) Update(ctx context.Context, question *entity.Question) error {
	return s.db.WithContext(ctx).Save(question).Error
}
