package dto

type UserResponse struct {
	ID             string `json:"user_id"`
	Fullname       string `json:"full_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	EducationLevel string `json:"education_level,omitempty"`
	Major          string `json:"major,omitempty"`
	Institution    string `json:"institution,omitempty"`
	GraduationYear int    `json:"graduation_year,omitempty"`
	IsPremium      bool   `json:"is_premium"`
	IsVerified     bool   `json:"is_verified"`
	CreatedAt      string `json:"created_at"`
}

type UsersUpdateRequest struct {
	Fullname       string `json:"full_name"`
	EducationLevel string `json:"education_level,omitempty"`
	Major          string `json:"major,omitempty"`
	Institution    string `json:"institution,omitempty"`
	GraduationYear int    `json:"graduation_year,omitempty"`
}
