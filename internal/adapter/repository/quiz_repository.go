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

func (q *quizRepository) GetRandomQuestionBySkillAndLevel(ctx context.Context, skillID string, level entity.LevelEnum, limit int) ([]entity.Question, error) {
	var questions []entity.Question
	err := q.db.WithContext(ctx).Where("skill_id = ? AND level = ?", skillID, level).Order("RANDOM()").Limit(limit).Find(&questions).Error
	return questions, err
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

func (q *quizRepository) UpdateQuizAnswer(ctx context.Context, quizSessionID string, quizAnswerID string, userAnswer string) error {
	return q.db.WithContext(ctx).Model(&entity.QuizAnswer{}).
		Where("id = ? AND quiz_session_id = ?", quizAnswerID, quizSessionID).
		Update("user_answer", userAnswer).Error
}

func (q *quizRepository) GetAnswerWithQuestions(ctx context.Context, quizSessionID string) ([]entity.QuizAnswer, error) {
	var quizAnswers []entity.QuizAnswer
	err := q.db.WithContext(ctx).Preload("Question").Preload("QuizSession").
		Where("quiz_session_id = ?", quizSessionID).Find(&quizAnswers).Error

	return quizAnswers, err
}

func (q *quizRepository) SubmitQuizTransaction(ctx context.Context, quizSessionID string, careerSessionID string, totalScore int, updatedSkill []entity.SelfAssessmentSkill, updatedAnswers []entity.QuizAnswer) error {
	tx := q.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Model(&entity.QuizSession{}).Where("id = ?", quizSessionID).Updates(map[string]interface{}{
		"status": "completed",
		"score":  totalScore,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, skill := range updatedSkill {
		err := tx.Model(&entity.SelfAssessmentSkill{}).Where("id = ?", skill.ID).Updates(map[string]interface{}{
			"user_final_level": skill.UserFinalLevel,
			"quiz_score":       skill.QuizScore,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, ans := range updatedAnswers {
		err := tx.Model(&entity.QuizAnswer{}).Where("id = ?", ans.ID).Update("is_correct", ans.IsCorrect).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Model(&entity.UserCareerSession{}).Where("id = ?", careerSessionID).Update("status", entity.StatusOnLearning).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (q *quizRepository) Delete(ctx context.Context, quizSessionID string) error {
	tx := q.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Where("quiz_session_id = ?", quizSessionID).Delete(&entity.QuizAnswer{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("id = ?", quizSessionID).Delete(&entity.QuizSession{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (q *quizRepository) GetQuizSessionStatus(ctx context.Context, careerSessionID string) (*entity.QuizSession, error) {
	var quizSession entity.QuizSession

	err := q.db.WithContext(ctx).Where("user_career_session_id = ? AND status = ?", careerSessionID, entity.StatusOnProcess).First(&quizSession).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &quizSession, nil
}
