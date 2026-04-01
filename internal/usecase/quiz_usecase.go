package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"strings"

	"github.com/google/uuid"
)

type QuizUsecase struct {
	quizRepo QuizRepository
}

func NewQuizUsecase(repo QuizRepository) *QuizUsecase {
	return &QuizUsecase{quizRepo: repo}
}

func isCorrectAnswer(question *entity.Question, answer string) bool {
	switch answer {
	case "A":
		return answer == question.Answer
	case "B":
		return answer == question.Answer
	case "C":
		return answer == question.Answer
	case "D":
		return answer == question.Answer
	default:
		return false
	}
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

func UserFinalLevel(userLevel entity.LevelEnum, isCorrect1 bool, isCorrect2 bool) entity.LevelEnum {
	totalCorrect := 0
	if isCorrect1 {
		totalCorrect++
	}
	if isCorrect2 {
		totalCorrect++
	}

	switch userLevel {
	case entity.LevelNoExperience:
		if totalCorrect == 2 {
			return entity.LevelBeginner
		}
		return entity.LevelNoExperience
	case entity.LevelBeginner:
		if totalCorrect == 2 {
			return entity.LevelIntermediate
		}
		if totalCorrect == 1 {
			return entity.LevelBeginner
		}
		return entity.LevelNoExperience
	case entity.LevelIntermediate:
		if totalCorrect == 2 {
			return entity.LevelExpert
		}
		if totalCorrect == 1 {
			return entity.LevelIntermediate
		}
		return entity.LevelBeginner
	case entity.LevelExpert:
		if totalCorrect >= 1 {
			return entity.LevelExpert
		}
		return entity.LevelIntermediate
	default:
		return entity.LevelNoExperience
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

	quizStatus, err := u.quizRepo.GetQuizSessionStatus(ctx, sessionUUID.String())
	if err != nil {
		return nil, errors.New("Gagal memeriksa status sesi kuis")
	}

	if quizStatus != nil {
		err := u.quizRepo.Delete(ctx, quizStatus.ID.String())
		if err != nil {
			return nil, errors.New("Gagal mereset sesi kuis")
		}
	}

	selfSkills, err := u.quizRepo.GetSelfAssessmentSkillsBySession(ctx, sessionUUID.String())
	if err != nil || len(selfSkills) == 0 {
		return nil, errors.New("Data self-assessment tidak ditemukan")
	}

	quizSessionUUID := uuid.New()
	quizSession := &entity.QuizSession{
		ID:                  quizSessionUUID,
		UserID:              userUUID,
		UserCareerSessionID: sessionUUID,
		Status:              entity.StatusOnProcess,
		Score:               0,
	}

	var quizAnswers []entity.QuizAnswer
	var questionsResponse []dto.QuizQuestionResponse

	for _, selfSkill := range selfSkills {
		level1, level2 := ShortingQuestionLevel(selfSkill.UserLevel)

		var selectedQuestion []entity.Question

		if level1 == level2 {
			quests, err := u.quizRepo.GetRandomQuestionBySkillAndLevel(ctx, selfSkill.SkillID.String(), level1, 2)
			if err != nil || len(quests) < 2 {
				return nil, errors.New("Bank soal tidak mencukupi untuk skill ini")
			}
			selectedQuestion = quests

		} else {
			quests1, err1 := u.quizRepo.GetRandomQuestionBySkillAndLevel(ctx, selfSkill.SkillID.String(), level1, 1)
			quests2, err2 := u.quizRepo.GetRandomQuestionBySkillAndLevel(ctx, selfSkill.SkillID.String(), level2, 1)
			if err1 != nil || err2 != nil || len(quests1) < 1 || len(quests2) < 1 {
				return nil, errors.New("Bank soal tidak mencukupi untuk skill ini")
			}
			selectedQuestion = append(selectedQuestion, quests1[0], quests2[0])
		}

		for _, q := range selectedQuestion {
			qa := entity.QuizAnswer{
				ID:         uuid.New(),
				QuestionID: q.ID,
				IsCorrect:  false,
			}
			quizAnswers = append(quizAnswers, qa)

			questionsResponse = append(questionsResponse, dto.QuizQuestionResponse{
				QuizAnswerID:    qa.ID.String(),
				QuestionID:      q.ID.String(),
				SkillID:         selfSkill.SkillID.String(),
				SkillName:       selfSkill.Skill.Name,
				QuestionContent: q.QuestionContent,
				OptionA:         q.OptionA,
				OptionB:         q.OptionB,
				OptionC:         q.OptionC,
				OptionD:         q.OptionD,
			})
		}

	}

	if err := u.quizRepo.CreateQuizTransaction(ctx, quizSession, quizAnswers); err != nil {
		return nil, err
	}

	return &dto.StartQuizResponse{
		QuizSessionID:   quizSessionUUID.String(),
		CareerSessionID: sessionUUID.String(),
		Questions:       questionsResponse,
	}, nil
}

func (u *QuizUsecase) UpdateAnswer(ctx context.Context, quizSessionID string, req dto.UpdateAnswerRequest) error {
	err := u.quizRepo.UpdateQuizAnswer(ctx, quizSessionID, req.QuizAnswerID, req.UserAnswer)
	if err != nil {
		return errors.New("Gagal menyimpan jawaban")
	}
	return nil
}

func (u *QuizUsecase) SubmitQuiz(ctx context.Context, userId string, quizSessionID string) (*dto.SubmitQuizResponse, error) {
	quizAnswers, err := u.quizRepo.GetAnswerWithQuestions(ctx, quizSessionID)
	if err != nil || len(quizAnswers) == 0 {
		return nil, errors.New("Gagal mengambil jawaban kuis")
	}

	careerSessionId := quizAnswers[0].QuizSession.UserCareerSessionID.String()
	selfAssessmentSkills, err := u.quizRepo.GetSelfAssessmentSkillsBySession(ctx, careerSessionId)
	if err != nil {
		return nil, errors.New("Gagal mengambil data self-assessment")
	}

	userLevel := make(map[string]entity.LevelEnum)
	for _, skill := range selfAssessmentSkills {
		userLevel[skill.SkillID.String()] = skill.UserLevel
	}

	answerSkill := make(map[string][]entity.QuizAnswer)
	for _, answer := range quizAnswers {
		skillID := answer.Question.SkillID.String()
		answerSkill[skillID] = append(answerSkill[skillID], answer)
	}

	totalScore := 0
	var skillsResult []dto.SkillFinalDetails
	var updatedSkills []entity.SelfAssessmentSkill
	var updatedAnswers []entity.QuizAnswer

	for _, ss := range selfAssessmentSkills {
		skillIDStr := ss.SkillID.String()
		answers := answerSkill[skillIDStr]

		if len(answers) != 2 {
			continue
		}

		isCorrect1 := strings.EqualFold(answers[0].UserAnswer, answers[0].Question.Answer)
		isCorrect2 := strings.EqualFold(answers[1].UserAnswer, answers[1].Question.Answer)

		answers[0].IsCorrect = isCorrect1
		answers[1].IsCorrect = isCorrect2
		updatedAnswers = append(updatedAnswers, answers[0], answers[1])

		skillScore := 0
		if isCorrect1 {
			totalScore += 10
			skillScore += 10
		}
		if isCorrect2 {
			totalScore += 10
			skillScore += 10
		}

		userLevelBefore := userLevel[skillIDStr]
		finalLevel := UserFinalLevel(userLevelBefore, isCorrect1, isCorrect2)

		skillsResult = append(skillsResult, dto.SkillFinalDetails{
			SkillID:    skillIDStr,
			SkillName:  ss.Skill.Name,
			UserLevel:  string(userLevelBefore),
			FinalLevel: string(finalLevel),
			SkillScore: skillScore,
		})

		ss.UserFinalLevel = finalLevel
		ss.QuizScore = skillScore
		updatedSkills = append(updatedSkills, ss)
	}

	err = u.quizRepo.SubmitQuizTransaction(ctx, quizSessionID, careerSessionId, totalScore, updatedSkills, updatedAnswers)
	if err != nil {
		return nil, errors.New("Gagal menyimpan hasil akhir kuis")
	}

	return &dto.SubmitQuizResponse{
		QuizSessionID: quizSessionID,
		TotalScore:    totalScore,
		SkillsResult:  skillsResult,
	}, nil
}
