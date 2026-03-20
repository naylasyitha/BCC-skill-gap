package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"
)

type CareerUsecase struct {
	careerRepository CareerRepository
}

func NewCareerUsecase(repo CareerRepository) *CareerUsecase {
	return &CareerUsecase{careerRepository: repo}
}

func (cu *CareerUsecase) GetAllCareer(ctx context.Context) ([]dto.CareerResponse, error) {
	careers, err := cu.careerRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var responses []dto.CareerResponse
	for _, career := range careers {
		responses = append(responses, dto.CareerResponse{
			ID:   career.ID.String(),
			Name: career.Name,
			Desc: career.Desc,
		})
	}
	return responses, nil
}

func (cu *CareerUsecase) GetCareerById(ctx context.Context, id string) (*dto.CareerSkillResponse, error) {
	careers, err := cu.careerRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	var skillResponses []dto.SkillsResponse
	for _, careerSkill := range careers.CareerSkills {
		skillResponses = append(skillResponses, dto.SkillsResponse{
			ID:            careerSkill.Skill.ID.String(),
			Name:          careerSkill.Skill.Name,
			Desc:          careerSkill.Skill.Desc,
			Priority:      careerSkill.Priority,
			RequiredLevel: string(careerSkill.RequiredLevel),
		})
	}

	return &dto.CareerSkillResponse{
		ID:     careers.ID.String(),
		Name:   careers.Name,
		Desc:   careers.Desc,
		Skills: skillResponses,
	}, nil
}

func (cu *CareerUsecase) CreateCareer(ctx context.Context, req dto.CareerCreateRequest) (*dto.CareerResponse, error) {
	career := &entity.Career{
		Name: req.Name,
		Desc: req.Desc,
	}
	err := cu.careerRepository.Save(ctx, career)

	if err != nil {
		return nil, err
	}
	return &dto.CareerResponse{
		ID:   career.ID.String(),
		Name: career.Name,
		Desc: career.Desc,
	}, nil
}

func (cu *CareerUsecase) UpdateCareer(ctx context.Context, id string, req dto.CareerEditRequest) (*dto.CareerResponse, error) {
	career, err := cu.careerRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		career.Name = req.Name
	}

	if req.Desc != "" {
		career.Desc = req.Desc
	}

	err = cu.careerRepository.Update(ctx, career)
	if err != nil {
		return nil, err
	}

	return &dto.CareerResponse{
		ID:   career.ID.String(),
		Name: career.Name,
		Desc: career.Desc,
	}, nil
}

func (cu *CareerUsecase) DeleteCareer(ctx context.Context, id string) error {

	_, err := cu.careerRepository.FindById(ctx, id)
	if err != nil {
		return errors.New("Karir tidak ditemukan")
	}
	err = cu.careerRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
