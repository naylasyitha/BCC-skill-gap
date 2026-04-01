package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"

	"github.com/google/uuid"
)

type SelfAssessmentUsecase struct {
	selfAssessmentRepository SelfAssessmentRepository
	careerSessionRepo        CareerSessionRepository
}

func NewSelfAssessmentUsecase(repo SelfAssessmentRepository, careerSessionRepo CareerSessionRepository) *SelfAssessmentUsecase {
	return &SelfAssessmentUsecase{selfAssessmentRepository: repo, careerSessionRepo: careerSessionRepo}
}

func (s *SelfAssessmentUsecase) ProcessSelfAssessment(ctx context.Context, careerSessionID string, req dto.SelfAssessmentRequest) (*dto.SelfAssessmentResponse, error) {

	careerSessionUUID, err := uuid.Parse(careerSessionID)
	if err != nil {
		return nil, errors.New("Career Session ID tidak valid")
	}

	careerSession, err := s.careerSessionRepo.FindById(ctx, careerSessionID)
	if err != nil {
		return nil, errors.New("Career Session tidak ditemukan")
	}

	if careerSession.Status != entity.StatusEnum("on_assessment") {
		return nil, errors.New("Status career session tidak sesuai untuk self assessment")
	}

	var skills []entity.SelfAssessmentSkill

	for _, skillReq := range req.Skills {
		skillUUID, err := uuid.Parse(skillReq.SkillID)
		if err != nil {
			return nil, errors.New("Skill ID tidak valid")
		}

		if skillReq.UserLevel == "" {
			skillReq.UserLevel = string(entity.LevelNoExperience)
		}

		skills = append(skills, entity.SelfAssessmentSkill{
			UserCareerSessionID: careerSessionUUID,
			SkillID:             skillUUID,
			UserLevel:           entity.LevelEnum(skillReq.UserLevel),
		})
	}

	if err := s.selfAssessmentRepository.CreateAssessmentSession(ctx, skills); err != nil {
		return nil, err
	}

	if err = s.selfAssessmentRepository.UpdateStatus(ctx, careerSessionID, entity.StatusOnQuiz); err != nil {
		return nil, err
	}

	response := &dto.SelfAssessmentResponse{
		UserCareerSessionID: careerSessionUUID.String(),
	}
	return response, nil
}
