package usecase

import (
	"context"
	"errors"
	"fmt"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidSessionStatus = errors.New("Status career session tidak sesuai untuk self assessment")
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
		return nil, errors.New("Format Career Session ID tidak valid")
	}

	careerSession, err := s.careerSessionRepo.FindById(ctx, careerSessionUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrCareerSessionNotFound
		}
		fmt.Println("Gagal mengambil data Career Session:", err.Error())
		return nil, ErrInternalServer
	}

	if careerSession.Status != entity.StatusEnum("on_assessment") {
		return nil, ErrInvalidSessionStatus
	}

	var skills []entity.SelfAssessmentSkill

	for _, skillReq := range req.Skills {
		skillUUID, err := uuid.Parse(skillReq.SkillID)
		if err != nil {
			return nil, errors.New("Format Skill ID tidak valid")
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
		fmt.Println("Gagal menyimpan Self Assessment ke database:", err.Error())
		return nil, ErrInternalServer
	}

	if err = s.selfAssessmentRepository.UpdateStatus(ctx, careerSessionUUID.String(), entity.StatusOnQuiz); err != nil {
		fmt.Println("Gagal update status Career Session ke on_quiz:", err.Error())
		return nil, ErrInternalServer
	}

	response := &dto.SelfAssessmentResponse{
		UserCareerSessionID: careerSessionUUID.String(),
	}
	return response, nil
}
