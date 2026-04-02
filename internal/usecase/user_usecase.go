package usecase

import (
	"context"
	"errors"
	"fmt"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("User tidak ditemukan")
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		fmt.Println("Gagal mengambil data user:", err.Error())
		return nil, ErrInternalServer
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
		fmt.Println("Gagal menyimpan update user:", err.Error())
		return nil, ErrInternalServer
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		fmt.Println("Gagal mengambil data user:", err.Error())
		return ErrInternalServer
	}

	err = u.authRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		fmt.Println("Gagal menghapus data:", err.Error())
		return ErrInternalServer
	}
	return nil
}
