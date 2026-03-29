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
}

func NewSelfAssessmentUsecase(repo SelfAssessmentRepository) *SelfAssessmentUsecase {
	return &SelfAssessmentUsecase{selfAssessmentRepository: repo}
}

func (s *SelfAssessmentUsecase) ProcessSelfAssessment(ctx context.Context, careerSessionID string, req dto.SelfAssessmentRequest) (dto.SelfAssessmentResponse, error) {
	var response dto.SelfAssessmentResponse

	careerSessionUUID, err := uuid.Parse(careerSessionID)
	if err != nil {
		return response, errors.New("Career Session ID tidak valid")
	}

	var skills []entity.SelfAssessmentSkill

	for _, skillReq := range req.Skills {
		skillUUID, err := uuid.Parse(skillReq.SkillID)
		if err != nil {
			return response, errors.New("Skill ID tidak valid")
		}

		skills = append(skills, entity.SelfAssessmentSkill{
			UserCareerSessionID: careerSessionUUID,
			SkillID:             skillUUID,
			UserLevel:           entity.LevelEnum(skillReq.UserLevel),
		})
	}

	if err := s.selfAssessmentRepository.CreateAssessmentSession(ctx, skills); err != nil {
		return response, err
	}

	// Langkah E: Kembalikan ID sesi yang berhasil dibuat ke dalam Response DTO
	response.UserCareerSessionID = careerSessionUUID.String()
	return response, nil
}
