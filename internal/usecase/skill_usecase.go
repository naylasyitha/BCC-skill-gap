package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"

	"github.com/google/uuid"
)

type SkillUsecase struct {
	skillRepository       SkillRepository
	careerSkillRepository CareerSkillRepository
}

func NewSkillUsecase(
	repo SkillRepository,
	repoCS CareerSkillRepository,
) *SkillUsecase {
	return &SkillUsecase{
		skillRepository:       repo,
		careerSkillRepository: repoCS,
	}
}

func (s *SkillUsecase) GetAllSkill(ctx context.Context) ([]dto.SkillResponse, error) {
	skills, err := s.skillRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.SkillResponse
	for _, skill := range skills {
		responses = append(responses, dto.SkillResponse{
			ID:   skill.ID.String(),
			Name: skill.Name,
			Desc: skill.Desc,
		})
	}

	return responses, nil
}

func (s *SkillUsecase) GetSkillById(ctx context.Context, id string) (*dto.SkillResponse, error) {
	skill, err := s.skillRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.SkillResponse{
		ID:   skill.ID.String(),
		Name: skill.Name,
		Desc: skill.Desc,
	}, nil
}

func (s *SkillUsecase) CreateSkill(ctx context.Context, req dto.SkillCreateRequest) (*dto.SkillResponse, error) {
	skill := &entity.Skill{
		Name: req.Name,
		Desc: req.Desc,
	}

	err := s.skillRepository.Save(ctx, skill)
	if err != nil {
		return nil, err
	}

	return &dto.SkillResponse{
		ID:   skill.ID.String(),
		Name: skill.Name,
		Desc: skill.Desc,
	}, nil

}

func (s *SkillUsecase) UpdateSkill(ctx context.Context, id string, req dto.SkillUpdateRequest) (*dto.SkillResponse, error) {
	skill, err := s.skillRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		skill.Name = req.Name
	}

	if req.Desc != "" {
		skill.Desc = req.Desc
	}

	err = s.skillRepository.Update(ctx, skill)
	if err != nil {
		return nil, err
	}

	return &dto.SkillResponse{
		ID:   skill.ID.String(),
		Name: skill.Name,
		Desc: skill.Desc,
	}, nil
}

func (s *SkillUsecase) DeleteSkill(ctx context.Context, id string) error {
	_, err := s.skillRepository.FindById(ctx, id)
	if err != nil {
		return errors.New("Skill tidak ditemukan")
	}

	err = s.skillRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SkillUsecase) CareerSkillAsign(ctx context.Context, req dto.CareerSkillCreateRequest) (*dto.CareerSkillAsignResponse, error) {
	careerSkill := &entity.CareerSkill{
		CareerID:      uuid.MustParse(req.CareerID),
		SkillID:       uuid.MustParse(req.SkillID),
		Priority:      req.Priority,
		RequiredLevel: entity.LevelEnum(req.RequiredLevel),
	}
	err := s.careerSkillRepository.Save(ctx, careerSkill)
	if err != nil {
		return nil, err
	}

	return &dto.CareerSkillAsignResponse{
		ID:            careerSkill.ID.String(),
		CareerID:      careerSkill.CareerID.String(),
		SkillID:       careerSkill.SkillID.String(),
		Priority:      careerSkill.Priority,
		RequiredLevel: string(careerSkill.RequiredLevel),
	}, nil
}

func (s *SkillUsecase) UpdateCareerSkill(ctx context.Context, id string, req dto.CareerSkillUpdateRequest) (*dto.CareerSkillAsignResponse, error) {
	careerSkill, err := s.careerSkillRepository.FindById(ctx, id)
	if err != nil {
		return nil, errors.New("Karir skill tidak ditemukan")
	}
	if req.Priority != 0 {
		careerSkill.Priority = req.Priority
	}

	if req.RequiredLevel != "" {
		careerSkill.RequiredLevel = entity.LevelEnum(req.RequiredLevel)
	}

	err = s.careerSkillRepository.Update(ctx, careerSkill)
	if err != nil {
		return nil, err
	}

	return &dto.CareerSkillAsignResponse{
		ID:            careerSkill.ID.String(),
		CareerID:      careerSkill.CareerID.String(),
		SkillID:       careerSkill.SkillID.String(),
		Priority:      careerSkill.Priority,
		RequiredLevel: string(careerSkill.RequiredLevel),
	}, nil
}

func (s *SkillUsecase) RemoveSkillFromCareer(ctx context.Context, id string) error {
	_, err := s.careerSkillRepository.FindById(ctx, id)
	if err != nil {
		return errors.New("Karir skill tidak ditemukan")
	}

	err = s.careerSkillRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
