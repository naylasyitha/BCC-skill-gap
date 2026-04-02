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

var (
	ErrSkillNotFound  = errors.New("Skill tidak ditemukan")
	ErrCareerNotFound = errors.New("Karir tidak ditemukan")
)

type CareerUsecase struct {
	careerRepository CareerRepository
	skillRepository  SkillRepository
}

func NewCareerUsecase(
	repo CareerRepository,
	sRepo SkillRepository,
) *CareerUsecase {
	return &CareerUsecase{
		careerRepository: repo,
		skillRepository:  sRepo,
	}
}

func (cu *CareerUsecase) GetAllCareer(ctx context.Context) ([]dto.CareerResponse, error) {
	careers, err := cu.careerRepository.FindAll(ctx)
	if err != nil {
		fmt.Println("Gagal mengambil data semua karir:", err.Error())
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCareerNotFound
		}
		fmt.Println("Gagal mengambil data karir by ID:", err.Error())
		return nil, ErrInternalServer
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

func (cu *CareerUsecase) CreateCareer(ctx context.Context, req dto.CareerCreateRequest) (*dto.CareerSkillResponse, error) {
	career := &entity.Career{
		Name: req.Name,
		Desc: req.Desc,
	}

	var careerSkills []entity.CareerSkill
	var skillResponses []dto.SkillsResponse

	for _, reqSkill := range req.Skills {
		skillUUID, err := uuid.Parse(reqSkill.SkillID)
		if err != nil {
			return nil, errors.New("Format ID Skill tidak valid")
		}

		skillData, err := cu.skillRepository.FindById(ctx, reqSkill.SkillID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrSkillNotFound
			}
			fmt.Println("Gagal validasi skill saat CreateCareer:", err.Error())
			return nil, ErrInternalServer
		}

		careerSkills = append(careerSkills, entity.CareerSkill{
			SkillID:       skillUUID,
			Priority:      reqSkill.Priority,
			RequiredLevel: entity.LevelEnum(reqSkill.RequiredLevel),
		})

		skillResponses = append(skillResponses, dto.SkillsResponse{
			ID:            skillData.ID.String(),
			Name:          skillData.Name,
			Desc:          skillData.Desc,
			Priority:      reqSkill.Priority,
			RequiredLevel: reqSkill.RequiredLevel,
		})
	}

	err := cu.careerRepository.CreateCareerSkill(ctx, career, careerSkills)
	if err != nil {
		fmt.Println("Gagal menyimpan relasi Career & Skill ke database:", err.Error())
		return nil, ErrInternalServer
	}

	return &dto.CareerSkillResponse{
		ID:     career.ID.String(),
		Name:   career.Name,
		Desc:   career.Desc,
		Skills: skillResponses,
	}, nil
}

func (cu *CareerUsecase) UpdateCareer(ctx context.Context, id string, req dto.CareerEditRequest) (*dto.CareerSkillResponse, error) {
	career, err := cu.careerRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCareerNotFound
		}
		fmt.Println("Gagal query karir saat UpdateCareer:", err.Error())
		return nil, ErrInternalServer
	}

	if req.Name != "" {
		career.Name = req.Name
	}

	if req.Desc != "" {
		career.Desc = req.Desc
	}

	var newSkills []entity.CareerSkill
	updateSkills := req.Skills != nil

	if updateSkills {
		for _, s := range req.Skills {
			newSkills = append(newSkills, entity.CareerSkill{
				SkillID:       uuid.MustParse(s.SkillID),
				Priority:      s.Priority,
				RequiredLevel: entity.LevelEnum(s.RequiredLevel),
			})
		}
	}

	err = cu.careerRepository.UpdateCareerWithSkills(ctx, career, newSkills, updateSkills)
	if err != nil {
		fmt.Println("Gagal update data karir ke database:", err.Error())
		return nil, ErrInternalServer
	}

	var skillResponses []dto.SkillsResponse
	if updateSkills {
		for _, s := range req.Skills {
			skillDetail, err := cu.skillRepository.FindById(ctx, s.SkillID)
			if err == nil {
				skillResponses = append(skillResponses, dto.SkillsResponse{
					ID:            s.SkillID,
					Name:          skillDetail.Name,
					Desc:          skillDetail.Desc,
					Priority:      s.Priority,
					RequiredLevel: s.RequiredLevel,
				})
			}
		}
	} else {
		//lupa menambahkan skill lama jika skill tidak di update
		for _, cs := range career.CareerSkills {
			skillResponses = append(skillResponses, dto.SkillsResponse{
				ID:            cs.Skill.ID.String(),
				Name:          cs.Skill.Name,
				Desc:          cs.Skill.Desc,
				Priority:      cs.Priority,
				RequiredLevel: string(cs.RequiredLevel),
			})
		}
	}

	return &dto.CareerSkillResponse{
		ID:     career.ID.String(),
		Name:   career.Name,
		Desc:   career.Desc,
		Skills: skillResponses,
	}, nil
}

func (cu *CareerUsecase) DeleteCareer(ctx context.Context, id string) error {
	_, err := cu.careerRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCareerNotFound
		}
		fmt.Println("Gagal query karir saat DeleteCareer:", err.Error())
		return ErrInternalServer
	}

	err = cu.careerRepository.Delete(ctx, id)
	if err != nil {
		fmt.Println("Gagal menghapus karir di database:", err.Error())
		return ErrInternalServer
	}
	return nil
}
