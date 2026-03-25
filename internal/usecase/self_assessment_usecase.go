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

func (s *SelfAssessmentUsecase) ProcessSelfAssessment(ctx context.Context, userID string, req dto.SelfAssessmentRequest) (dto.SelfAssessmentResponse, error) {
	var response dto.SelfAssessmentResponse

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return response, errors.New("User ID tidak valid")
	}

	careerUUID, err := uuid.Parse(req.CareerID)
	if err != nil {
		return response, errors.New("Career ID tidak valid")
	}
	session := &entity.UserCareerSession{
		UserID:   userUUID,
		CareerID: careerUUID,
		Status:   entity.StatusOnProcess,
	}

	var skills []entity.SelfAssessmentSkill

	for _, skillReq := range req.Skills {
		skillUUID, err := uuid.Parse(skillReq.SkillID)
		if err != nil {
			return response, errors.New("Skill ID tidak valid")
		}
		skills = append(skills, entity.SelfAssessmentSkill{
			SkillID:   skillUUID,
			UserLevel: entity.LevelEnum(skillReq.UserLevel),
		})
	}

	if err := s.selfAssessmentRepository.CreateAssessmentSession(ctx, session, skills); err != nil {
		return response, err
	}

	// Langkah E: Kembalikan ID sesi yang berhasil dibuat ke dalam Response DTO
	response.UserCareerSessionID = session.ID.String()
	return response, nil
}
