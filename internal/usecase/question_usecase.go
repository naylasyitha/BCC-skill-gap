package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"

	"github.com/google/uuid"
)

type QuestionUsecase struct {
	questionRepo QuestionRepository
	skillRepo    SkillRepository
}

func NewQuestionUsecase(questRepo QuestionRepository, skill SkillRepository) *QuestionUsecase {
	return &QuestionUsecase{
		questionRepo: questRepo,
		skillRepo:    skill,
	}
}

func (u *QuestionUsecase) CreateQuestion(ctx context.Context, req dto.QuestionCreateRequest) (*dto.QuestionResponse, error) {
	skillUUID, err := uuid.Parse(req.SkillID)
	if err != nil {
		return nil, errors.New("Skill ID tidak valid")
	}

	_, err = u.skillRepo.FindById(ctx, req.SkillID)
	if err != nil {
		return nil, errors.New("Skill tidak ditemukan")
	}

	newQuestion := &entity.Question{
		ID:              uuid.New(),
		SkillID:         skillUUID,
		Level:           entity.LevelEnum(req.Level),
		QuestionContent: req.QuestionContent,
		OptionA:         req.OptionA,
		OptionB:         req.OptionB,
		OptionC:         req.OptionC,
		OptionD:         req.OptionD,
		Answer:          req.Answer,
		Explanation:     req.Explanation,
	}

	err = u.questionRepo.Create(ctx, newQuestion)
	if err != nil {
		return nil, err
	}

	return &dto.QuestionResponse{
		ID:              newQuestion.ID.String(),
		SkillID:         newQuestion.SkillID.String(),
		Level:           string(newQuestion.Level),
		QuestionContent: newQuestion.QuestionContent,
		OptionA:         newQuestion.OptionA,
		OptionB:         newQuestion.OptionB,
		OptionC:         newQuestion.OptionC,
		OptionD:         newQuestion.OptionD,
		Answer:          newQuestion.Answer,
		Explanation:     newQuestion.Explanation,
	}, nil
}

func (u *QuestionUsecase) GetAllQuestion(ctx context.Context) ([]dto.QuestionResponse, error) {
	questions, err := u.questionRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.QuestionResponse
	for _, question := range questions {
		responses = append(responses, dto.QuestionResponse{
			ID:              question.ID.String(),
			SkillID:         question.SkillID.String(),
			Level:           string(question.Level),
			QuestionContent: question.QuestionContent,
			OptionA:         question.OptionA,
			OptionB:         question.OptionB,
			OptionC:         question.OptionC,
			OptionD:         question.OptionD,
			Answer:          question.Answer,
			Explanation:     question.Explanation,
		})
	}
	return responses, nil
}

func (u *QuestionUsecase) GetQuestionById(ctx context.Context, id string) (*dto.QuestionResponse, error) {
	question, err := u.questionRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.QuestionResponse{
		ID:              question.ID.String(),
		SkillID:         question.SkillID.String(),
		Level:           string(question.Level),
		QuestionContent: question.QuestionContent,
		OptionA:         question.OptionA,
		OptionB:         question.OptionB,
		OptionC:         question.OptionC,
		OptionD:         question.OptionD,
		Answer:          question.Answer,
		Explanation:     question.Explanation,
	}, nil
}

func (u *QuestionUsecase) UpdateQuestion(ctx context.Context, id string, req dto.QuestionUpdateRequest) (*dto.QuestionResponse, error) {
	question, err := u.questionRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Level != "" {
		question.Level = entity.LevelEnum(req.Level)
	}
	if req.QuestionContent != "" {
		question.QuestionContent = req.QuestionContent
	}
	if req.OptionA != "" {
		question.OptionA = req.OptionA
	}
	if req.OptionB != "" {
		question.OptionB = req.OptionB
	}
	if req.OptionC != "" {
		question.OptionC = req.OptionC
	}
	if req.OptionD != "" {
		question.OptionD = req.OptionD
	}
	if req.Answer != "" {
		question.Answer = req.Answer
	}
	if req.Explanation != "" {
		question.Explanation = req.Explanation
	}

	err = u.questionRepo.Update(ctx, question)
	if err != nil {
		return nil, err
	}

	return &dto.QuestionResponse{
		ID:              question.ID.String(),
		SkillID:         question.SkillID.String(),
		Level:           string(question.Level),
		QuestionContent: question.QuestionContent,
		OptionA:         question.OptionA,
		OptionB:         question.OptionB,
		OptionC:         question.OptionC,
		OptionD:         question.OptionD,
		Answer:          question.Answer,
		Explanation:     question.Explanation,
	}, nil
}

func (u *QuestionUsecase) DeleteQuestion(ctx context.Context, id string) error {
	_, err := u.questionRepo.FindById(ctx, id)
	if err != nil {
		return errors.New("Question tidak ditemukan")
	}

	err = u.questionRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
