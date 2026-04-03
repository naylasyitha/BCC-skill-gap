package usecase

import (
	"context"
	"errors"
	"fmt"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrCareerSessionNotFound = errors.New("Career Session tidak ditemukan")
	ErrCareerLimitReached    = errors.New("Karir sudah mencapai batas memilih 2 karir. Silahkan upgrade menjadi premium.")
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

func getLevelString(levelInt int) string {
	switch levelInt {
	case 0:
		return "no_experience"
	case 1:
		return "beginner"
	case 2:
		return "intermediate"
	case 3:
		return "expert"
	default:
		return ""
	}
}

func getLevelInt(levelString string) int {
	switch levelString {
	case "no_experience":
		return 0
	case "beginner":
		return 1
	case "intermediate":
		return 2
	case "expert":
		return 3
	default:
		return -1
	}
}

func (cs *CareerSessionUsecase) CreateCareerSession(ctx context.Context, userID string, req dto.CareerSessionCreateRequest) (*dto.CareerSessionResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("Format User ID tidak valid")
	}

	user, err := cs.authRepo.FindByID(ctx, userUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			return nil, ErrUserNotFound
		}
		fmt.Println("Gagal mengambil data user saat CreateCareerSession:", err.Error())
		return nil, ErrInternalServer
	}

	if !user.IsPremium {
		count, err := cs.careerSessionRepo.CountByUserID(ctx, user.ID.String())
		if err != nil {
			fmt.Println("Gagal menghitung jumlah CareerSession:", err.Error())
			return nil, ErrInternalServer
		}
		if count >= 2 {
			return nil, ErrCareerLimitReached
		}
	}

	careerUUID, err := uuid.Parse(req.CareerID)
	if err != nil {
		return nil, errors.New("Format ID karir tidak valid")
	}

	_, err = cs.careerRepo.FindById(ctx, careerUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrCareerNotFound
		}
		fmt.Println("Gagal mencari Karir saat CreateCareerSession:", err.Error())
		return nil, ErrInternalServer
	}

	sessions := &entity.UserCareerSession{
		UserID:    user.ID,
		CareerID:  careerUUID,
		Status:    entity.StatusOnAssessment,
		StartedAt: time.Now(),
	}

	err = cs.careerSessionRepo.Create(ctx, sessions)
	if err != nil {
		fmt.Println("Gagal menyimpan Career Session:", err.Error())
		return nil, ErrInternalServer
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
	careerSessionUUID, err := uuid.Parse(careerSessionID)
	if err != nil {
		return nil, errors.New("Format Career Session ID tidak valid")
	}

	careerSession, err := cs.careerSessionRepo.FindById(ctx, careerSessionUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, ErrCareerSessionNotFound
		}
		fmt.Println("Gagal mencari Career Session Detail:", err.Error())
		return nil, ErrInternalServer
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
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("Format User ID tidak valid")
	}
	careerSessions, err := cs.careerSessionRepo.GetAllCareerSession(ctx, userUUID.String())
	if err != nil {
		fmt.Println("Gagal mengambil daftar Career Session:", err.Error())
		return nil, ErrInternalServer
	}

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
	careerSessionUUID, err := uuid.Parse(careerSessionID)
	if err != nil {
		return nil, errors.New("Format Career Session ID tidak valid")
	}

	assessments, err := cs.careerSessionRepo.GetAnalyticsData(ctx, careerSessionUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Data analitik tidak ditemukan untuk sesi ini")
		}
		return nil, ErrInternalServer
	}

	careerSkills, err := cs.careerSessionRepo.GetRequiredLevel(ctx, careerSessionUUID.String())
	if err != nil {
		return nil, ErrInternalServer
	}

	// Buat struct bantuan untuk map agar bisa menyimpan dua data sekaligus
	type skillReq struct {
		Level    string
		Priority int
	}

	industryLevel := make(map[string]skillReq)
	for _, required := range careerSkills {
		industryLevel[required.SkillID.String()] = skillReq{
			Level:    string(required.RequiredLevel),
			Priority: required.Priority, // Ambil priority dari data career_skill
		}
	}

	var totalScore int
	result := []dto.SkillAnalyticResponse{}
	var recommendations []dto.SkillRecommendation

	for _, a := range assessments {
		skillID := a.SkillID.String()
		skillPercentage := (float64(a.QuizScore) / 20.0) * 100.0
		totalScore += a.QuizScore

		// Ambil data requirement dari map
		reqData := industryLevel[skillID]
		requiredLevel := reqData.Level
		priority := reqData.Priority

		if requiredLevel == "" {
			requiredLevel = "beginner"
		}

		userLevelInt := getLevelInt(string(a.UserFinalLevel))
		reqLevelInt := getLevelInt(requiredLevel)
		gap := reqLevelInt - userLevelInt

		status := "Kemampuan anda masih kurang, Perlu ditingkatkan kembali!"
		var suggestion []string

		if gap <= 0 {
			status = "Selamat! kemampuan anda sudah sesuai dengan standar industri"
		} else {
			for i := userLevelInt + 1; i <= reqLevelInt; i++ {
				suggestion = append(suggestion, getLevelString(i))
			}

			// Masukkan data rekomendasi lengkap dengan Priority
			recommendations = append(recommendations, dto.SkillRecommendation{
				SkillID:      a.SkillID.String(),
				SkillName:    a.Skill.Name,
				CurrentLevel: string(a.UserFinalLevel),
				TargetLevel:  requiredLevel,
				Priority:     priority, // Data ini yang akan digunakan FE untuk sorting
			})
		}

		result = append(result, dto.SkillAnalyticResponse{
			SkillID:         a.SkillID.String(),
			SkillName:       a.Skill.Name,
			UserLevel:       string(a.UserLevel),
			FinalLevel:      string(a.UserFinalLevel),
			RequiredLevel:   requiredLevel,
			SkillScore:      int(skillPercentage),
			GapLevel:        gap,
			Status:          status,
			SuggestionLevel: suggestion,
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
		Recommendations: recommendations,
	}, nil
}
