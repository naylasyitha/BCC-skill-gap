package usecase

import (
	"context"
	"errors"
	"fmt"
	"project-bcc/dto"
	"project-bcc/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
		fmt.Println("Gagal mengambil semua data skill:", err.Error())
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSkillNotFound
		}
		fmt.Println("Gagal mengambil data skill by ID:", err.Error())
		return nil, ErrInternalServer
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
		fmt.Println("Gagal menyimpan skill baru:", err.Error())
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSkillNotFound
		}
		fmt.Println("Gagal query skill saat update:", err.Error())
		return nil, ErrInternalServer
	}

	if req.Name != "" {
		skill.Name = req.Name
	}

	if req.Desc != "" {
		skill.Desc = req.Desc
	}

	err = s.skillRepository.Update(ctx, skill)
	if err != nil {
		fmt.Println("Gagal menyimpan update skill:", err.Error())
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSkillNotFound
		}
		fmt.Println("Gagal query skill saat delete:", err.Error())
		return ErrInternalServer
	}

	err = s.skillRepository.Delete(ctx, id)
	if err != nil {
		fmt.Println("Gagal menghapus skill di database:", err.Error())
		return ErrInternalServer
	}
	return nil
}

func (s *SkillUsecase) CareerSkillAsign(ctx context.Context, req dto.CareerSkillCreateRequest) (*dto.CareerSkillAsignResponse, error) {
	careerUUID, err1 := uuid.Parse(req.CareerID)
	skillUUID, err2 := uuid.Parse(req.SkillID)

	if err1 != nil || err2 != nil {
		return nil, errors.New("Format ID Karir atau Skill tidak valid")
	}

	careerSkill := &entity.CareerSkill{
		CareerID:      careerUUID,
		SkillID:       skillUUID,
		Priority:      req.Priority,
		RequiredLevel: entity.LevelEnum(req.RequiredLevel),
	}
	err := s.careerSkillRepository.Save(ctx, careerSkill)
	if err != nil {
		fmt.Println("Gagal menyimpan relasi CareerSkill:", err.Error())
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Relasi karir dan skill tidak ditemukan")
		}
		fmt.Println("Gagal query CareerSkill saat update:", err.Error())
		return nil, ErrInternalServer
	}
	if req.Priority != 0 {
		careerSkill.Priority = req.Priority
	}

	if req.RequiredLevel != "" {
		careerSkill.RequiredLevel = entity.LevelEnum(req.RequiredLevel)
	}

	err = s.careerSkillRepository.Update(ctx, careerSkill)
	if err != nil {
		fmt.Println("Gagal menyimpan update CareerSkill:", err.Error())
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Relasi karir dan skill tidak ditemukan")
		}
		fmt.Println("Gagal query CareerSkill saat delete:", err.Error())
		return ErrInternalServer
	}

	err = s.careerSkillRepository.Delete(ctx, id)
	if err != nil {
		fmt.Println("Gagal menghapus relasi CareerSkill:", err.Error())
		return ErrInternalServer
	}
	return nil
}
