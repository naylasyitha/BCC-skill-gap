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
	authRepo          AuthRepository
}

func NewCareerSessionUsecase(
	csRepo CareerSessionRepository,
	cRepo CareerRepository,
	aRepo AuthRepository,
) *CareerSessionUsecase {
	return &CareerSessionUsecase{
		careerSessionRepo: csRepo,
		careerRepo:        cRepo,
		authRepo:          aRepo,
	}
}

func (cs *CareerSessionUsecase) CreateCareerSession(ctx context.Context, userID string, req dto.CareerSessionCreateRequest) (*dto.CareerSessionResponse, error) {
	user, err := cs.authRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("User tidak ditemukan")
	}

	if !user.IsPremium {
		count, err := cs.careerSessionRepo.CountByUserID(ctx, user.ID.String())
		if err != nil {
			return nil, errors.New("Gagal menghitung batas karir")
		}
		if count == 2 {
			return nil, errors.New("Gagal menambah karir, karir sudah mencapai batas memilih 2 karir. Silahkan upgrade menjadi premium untuk memilih karir lebih banyak.")
		}
	}

	_, err = cs.careerRepo.FindById(ctx, req.CareerID)
	if err != nil {
		return nil, errors.New("Karir tidak ditemukan")
	}
	userUUID := user.ID

	careerUUID, err := uuid.Parse(req.CareerID)
	if err != nil {
		return nil, errors.New("Format ID karir tidak valid")
	}

	sessions := &entity.UserCareerSession{
		UserID:    userUUID,
		CareerID:  careerUUID,
		Status:    entity.StatusOnAssessment,
		StartedAt: time.Now(),
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

func (cs *CareerSessionUsecase) GetAllCareerSession(ctx context.Context, userID string) ([]dto.CareerSessionListResponse, error) {
	careerSessions, err := cs.careerSessionRepo.GetAllCareerSession(ctx, userID)
	if err != nil {
		return nil, errors.New("Career Session tidak ditemukan")
	}

	//agar return array kosong bukan nil jika careerSession tidak ada
	responses := []dto.CareerSessionListResponse{}
	for _, careerSession := range careerSessions {
		responses = append(responses, dto.CareerSessionListResponse{
			CareerSessionID: careerSession.ID.String(),
			CareerID:        careerSession.CareerID.String(),
			CareerName:      careerSession.Career.Name,
			Status:          string(careerSession.Status),
		})
	}
	return responses, nil
}

func (cs *CareerSessionUsecase) GetDashboardAnalytics(ctx context.Context, careerSessionID string) (*dto.CareerAnalyticResponse, error) {
	assessments, err := cs.careerSessionRepo.GetAnalyticsData(ctx, careerSessionID)
	if err != nil {
		return nil, errors.New("Gagal mengambil data analitik")
	}

	var totalScore int
	result := []dto.SkillAnalyticResponse{}

	for _, a := range assessments {
		skillPercentage := (float64(a.QuizScore) / 20.0) * 100.0
		totalScore += a.QuizScore

		result = append(result, dto.SkillAnalyticResponse{
			SkillID:    a.SkillID.String(),
			SkillName:  a.Skill.Name,
			UserLevel:  string(a.UserLevel),
			FinalLevel: string(a.UserFinalLevel),
			SkillScore: int(skillPercentage),
		})
	}

	totalMaxScore := len(assessments) * 20
	totalPercentage := 0.0
	if totalMaxScore > 0 {
		totalPercentage = (float64(totalScore) / float64(totalMaxScore)) * 100.0
	}

	return &dto.CareerAnalyticResponse{
		CareerSessionID: careerSessionID,
		TotalScore:      int(totalPercentage),
		SkillsResult:    result,
	}, nil
}
