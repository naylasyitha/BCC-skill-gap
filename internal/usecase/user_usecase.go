package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"time"
)

type UserUsecase struct {
	authRepo AuthRepository
}

func NewUserUsecase(authRepo AuthRepository) *UserUsecase {
	return &UserUsecase{
		authRepo: authRepo,
	}
}

func (u *UserUsecase) GetUserData(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := u.authRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:             user.ID.String(),
		Fullname:       user.FullName,
		Email:          user.Email,
		Role:           string(user.Role),
		EducationLevel: string(user.EducationLevel),
		Major:          user.Major,
		Institution:    user.Institution,
		GraduationYear: user.GraduationYear,
		IsPremium:      user.IsPremium,
		IsVerified:     user.IsVerified,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, req dto.UsersUpdateRequest) (*dto.UserResponse, error) {
	user, err := u.authRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	if req.Fullname != "" {
		user.FullName = req.Fullname
	}

	if req.EducationLevel != "" {
		user.EducationLevel = entity.EduLevel(req.EducationLevel)
	}

	if req.Major != "" {
		user.Major = req.Major
	}

	if req.Institution != "" {
		user.Institution = req.Institution
	}

	if req.GraduationYear != 0 {
		user.GraduationYear = req.GraduationYear
	}

	err = u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.New("Gagal memperbarui data user")
	}

	return &dto.UserResponse{
		ID:             user.ID.String(),
		Fullname:       user.FullName,
		Email:          user.Email,
		Role:           string(user.Role),
		EducationLevel: string(user.EducationLevel),
		Major:          user.Major,
		Institution:    user.Institution,
		GraduationYear: user.GraduationYear,
		IsPremium:      user.IsPremium,
		IsVerified:     user.IsVerified,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
	}, nil
}
func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	_, err := u.authRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	err = u.authRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
