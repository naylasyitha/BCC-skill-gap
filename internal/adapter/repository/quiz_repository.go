package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) usecase.QuizRepository {
	return &quizRepository{db}
}

func (q *quizRepository) GetSelfAssessmentSkillsBySession(ctx context.Context, sessionID string) ([]entity.SelfAssessmentSkill, error) {
	var skills []entity.SelfAssessmentSkill
	err := q.db.WithContext(ctx).Where("user_career_session_id = ?", sessionID).Preload("Skill").Find(&skills).Error
	return skills, err
}

func (q *quizRepository) GetRandomQuestionBySkillAndLevel(ctx context.Context, skillID string, level entity.LevelEnum) (*entity.Question, error) {
	var question entity.Question
	err := q.db.WithContext(ctx).Where("skill_id = ? AND level = ?", skillID, level).Order("RANDOM()").First(&question).Error
	return &question, err
}

func (q *quizRepository) CreateQuizTransaction(ctx context.Context, quizSession *entity.QuizSession, quizAnswers []entity.QuizAnswer) error {
	tx := q.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(quizSession).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range quizAnswers {
		quizAnswers[i].QuizSessionID = quizSession.ID
	}

	if err := tx.Create(&quizAnswers).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
