package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"time"

	"github.com/google/uuid"
)

type CareerSessionUsecase struct {
	careerSessionRepo CareerSessionRepository
	careerRepo        CareerRepository
}

func NewCareerSessionUsecase(
	csRepo CareerSessionRepository,
	cRepo CareerRepository,
) *CareerSessionUsecase {
	return &CareerSessionUsecase{
		careerSessionRepo: csRepo,
		careerRepo:        cRepo,
	}
}

func (cs *CareerSessionUsecase) CreateCareerSession(ctx context.Context, userID string, req dto.CareerSessionCreateRequest) (*dto.CareerSessionResponse, error) {

	_, err := cs.careerRepo.FindById(ctx, req.CareerID)
	if err != nil {
		return nil, errors.New("Karir tidak ditemukan")
	}

	sessions := &entity.UserCareerSession{
		UserID:   uuid.MustParse(userID),
		CareerID: uuid.MustParse(req.CareerID),
		Status:   entity.StatusOnAssessment,
	}

	err = cs.careerSessionRepo.Create(ctx, sessions)
	if err != nil {
		return nil, err
	}

	completedAt := ""
	if sessions.CompletedAt != nil {
		completedAt = sessions.CompletedAt.Format(time.RFC3339)
	}

	return &dto.CareerSessionResponse{
		ID:          sessions.ID.String(),
		UserID:      sessions.UserID.String(),
		CareerID:    sessions.CareerID.String(),
		Status:      string(sessions.Status),
		StartedAt:   sessions.StartedAt.Format(time.RFC3339),
		CompletedAt: completedAt,
	}, nil
}

func (cs *CareerSessionUsecase) GetCareerSession(ctx context.Context, careerSessionID string) (*dto.CareerSessionDetailResponse, error) {
	careerSession, err := cs.careerSessionRepo.FindById(ctx, careerSessionID)
	if err != nil {
		return nil, errors.New("Career Session tidak ditemukan")
	}

	completedAt := ""
	if careerSession.CompletedAt != nil {
		completedAt = careerSession.CompletedAt.Format(time.RFC3339)
	}

	return &dto.CareerSessionDetailResponse{
		ID:          careerSession.ID.String(),
		UserID:      careerSession.UserID.String(),
		Fullname:    careerSession.User.FullName,
		CareerID:    careerSession.CareerID.String(),
		CareerName:  careerSession.Career.Name,
		Status:      string(careerSession.Status),
		StartedAt:   careerSession.StartedAt.Format(time.RFC3339),
		CompletedAt: completedAt,
	}, nil
}
