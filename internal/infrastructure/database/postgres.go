package database

import (
	"fmt"
	"log"
	"os"
	"project-bcc/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s timezone=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_TZ"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Koneksi database gagal", err)
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Career{},
		&entity.Skill{},
		&entity.CareerSkill{},
		&entity.UserCareerSession{},
		&entity.SelfAssessmentSkill{},
		&entity.Question{},
		&entity.QuizSession{},
		&entity.QuizAnswer{},
		&entity.Material{},
		&entity.LearningPathProgress{},
	)

	if err != nil {
		log.Fatal("Gagal membuat migrasi database")
	}
	log.Println("Koneksi database berhasil")

	return db
}
