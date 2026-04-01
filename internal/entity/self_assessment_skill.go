package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SelfAssessmentSkill struct {
	ID                  uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserCareerSessionID uuid.UUID `gorm:"type:uuid;not null"`
	UserCareerSession   UserCareerSession
	SkillID             uuid.UUID `gorm:"type:uuid;not null"`
	Skill               Skill
	UserLevel           LevelEnum `gorm:"type:varchar(50)"`
	UserFinalLevel      LevelEnum `gorm:"type:varchar(50)"`
	QuizScore           int       `gorm:"default:0"`
}

func (s *SelfAssessmentSkill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
