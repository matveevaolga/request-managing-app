package main

import "time"

type userSeed struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type typeSeed struct {
	Name string `json:"name"`
}

type applicationSeed struct {
	FullName         string     `json:"full_name"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone"`
	OrganisationName string     `json:"organisation_name"`
	OrganisationURL  *string    `json:"organisation_url"`
	ProjectName      string     `json:"project_name"`
	TypeName         string     `json:"type_name"`
	ExpectedResults  string     `json:"expected_results"`
	IsPayed          bool       `json:"is_payed"`
	Status           string     `json:"status"`
	RejectedReason   *string    `json:"rejected_reason"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}
