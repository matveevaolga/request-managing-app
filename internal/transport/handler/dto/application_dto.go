package dto

import "time"

type CreateApplicationRequest struct {
	FullName              string  `json:"fullName" validate:"required,max=150"`
	Email                 string  `json:"email" validate:"required,email,max=150"`
	Phone                 *string `json:"phone,omitempty" validate:"omitempty"`
	OrganisationName      string  `json:"organisationName" validate:"required,max=150"`
	OrganisationURL       *string `json:"organisationUrl,omitempty" validate:"omitempty,url"`
	ProjectName           string  `json:"projectName" validate:"required,max=150"`
	TypeID                int64   `json:"typeId" validate:"required,gt=0"`
	ExpectedResults       string  `json:"expectedResults" validate:"required,max=1500"`
	IsPayed               bool    `json:"isPayed"`
	AdditionalInformation *string `json:"additionalInformation,omitempty" validate:"omitempty,max=1500"`
}

type ApplicationResponse struct {
	ApplicationID         int64   `json:"applicationId"`
	FullName              string  `json:"fullName"`
	Email                 string  `json:"email"`
	Phone                 *string `json:"phone,omitempty"`
	OrganisationName      string  `json:"organisationName"`
	OrganisationURL       *string `json:"organisationUrl,omitempty"`
	ProjectName           string  `json:"projectName"`
	TypeName              string  `json:"typeName"`
	ExpectedResults       string  `json:"expectedResults"`
	IsPayed               bool    `json:"isPayed"`
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
	Status                string  `json:"status"`
}

type ApplicationPreviewResponse struct {
	ExternalApplicationID int64     `json:"externalApplicationId"`
	ProjectName           string    `json:"projectName"`
	TypeName              string    `json:"typeName"`
	Initiator             string    `json:"initiator"`
	OrganisationName      string    `json:"organisationName"`
	DateUpdated           time.Time `json:"dateUpdated"`
	Status                string    `json:"status"`
	RejectionMessage      *string   `json:"rejectionMessage,omitempty"`
}

type ApplicationListResponse struct {
	Count        int                          `json:"count"`
	Applications []ApplicationPreviewResponse `json:"applications"`
}

type RejectRequest struct {
	Reason string `json:"reason" validate:"required,max=750"`
}
