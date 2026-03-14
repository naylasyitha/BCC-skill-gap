package entity

type EduLevel string
type Role string
type LevelEnum string
type StatusEnum string

const (
	EduSMA EduLevel = "sma"
	EduSMK EduLevel = "smk"
	EduD3  EduLevel = "d3"
	EduD4  EduLevel = "d4"
	EduS1  EduLevel = "s1"

	RoleUser  Role = "user"
	RoleAdmin Role = "admin"

	LevelBeginner     LevelEnum = "beginner"
	LevelIntermediate LevelEnum = "intermediate"
	LevelExpert       LevelEnum = "expert"

	StatusComplete   StatusEnum = "complete"
	StatusOnProcess  StatusEnum = "on_process"
	StatusNotStarted StatusEnum = "not_started"
)
