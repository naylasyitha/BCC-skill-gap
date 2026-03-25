package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"

	"github.com/google/uuid"
)

type QuizUsecase struct {
	quizRepo QuizRepository
}

func NewQuizUsecase(repo QuizRepository) *QuizUsecase {
	return &QuizUsecase{quizRepo: repo}
}

func ShortingQuestionLevel(userLevel entity.LevelEnum) (entity.LevelEnum, entity.LevelEnum) {
	switch userLevel {
	case entity.LevelNoExperience:
		return entity.LevelBeginner, entity.LevelBeginner
	case entity.LevelBeginner:
		return entity.LevelBeginner, entity.LevelIntermediate
	case entity.LevelIntermediate:
		return entity.LevelIntermediate, entity.LevelExpert
	case entity.LevelExpert:
		return entity.LevelExpert, entity.LevelExpert
	default:
		return entity.LevelBeginner, entity.LevelBeginner
	}
}

func (u *QuizUsecase) StartQuiz(ctx context.Context, userID string, careerSessionID string) (*dto.StartQuizResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("User ID tidak valid")
	}

	sessionUUID, err := uuid.Parse(careerSessionID)
	if err != nil {
		return nil, errors.New("Career Session ID tidak valid")
	}

	selfSkills, err := u.quizRepo.GetSelfAssessmentSkillsBySession(ctx, sessionUUID.String())
	if err != nil || len(selfSkills) == 0 {
		return nil, errors.New("Data self-assessment tidak ditemukan")
	}

	quizSession := &entity.QuizSession{
		UserID:              userUUID,
		UserCareerSessionID: sessionUUID,
		Status:              entity.StatusOnProcess,
		Score:               0,
	}

	var quizAnswers []entity.QuizAnswer
	var questionsResponse []dto.QuizQuestionResponse

	for _, selfSkill := range selfSkills {
		level1, level2 := ShortingQuestionLevel(selfSkill.UserLevel)

		q1, err1 := u.quizRepo.GetRandomQuestionBySkillAndLevel(ctx, selfSkill.SkillID.String(), level1)
		q2, err2 := u.quizRepo.GetRandomQuestionBySkillAndLevel(ctx, selfSkill.SkillID.String(), level2)

		if err1 != nil || err2 != nil {
			return nil, errors.New("Bank soal tidak mencukupi untuk skill ini")
		}

		qa1 := entity.QuizAnswer{
			ID:         uuid.New(),
			QuestionID: q1.ID,
			IsCorrect:  false,
		}
		qa2 := entity.QuizAnswer{
			ID:         uuid.New(),
			QuestionID: q2.ID,
			IsCorrect:  false,
		}

		quizAnswers = append(quizAnswers, qa1, qa2)

		questionsResponse = append(questionsResponse, dto.QuizQuestionResponse{
			QuizAnswerID:    qa1.ID.String(),
			QuestionID:      q1.ID.String(),
			SkillID:         selfSkill.SkillID.String(),
			QuestionContent: q1.QuestionContent,
			OptionA:         q1.OptionA,
			OptionB:         q1.OptionB,
			OptionC:         q1.OptionC,
			OptionD:         q1.OptionD,
		})

		questionsResponse = append(questionsResponse, dto.QuizQuestionResponse{
			QuizAnswerID:    qa2.ID.String(),
			QuestionID:      q2.ID.String(),
			SkillID:         selfSkill.SkillID.String(),
			QuestionContent: q2.QuestionContent,
			OptionA:         q2.OptionA,
			OptionB:         q2.OptionB,
			OptionC:         q2.OptionC,
			OptionD:         q2.OptionD,
		})
	}

	if err := u.quizRepo.CreateQuizTransaction(ctx, quizSession, quizAnswers); err != nil {
		return nil, err
	}

	return &dto.StartQuizResponse{
		QuizSessionID: quizSession.ID.String(),
		Questions:     questionsResponse,
	}, nil
}
