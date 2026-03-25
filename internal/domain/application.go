package domain

import (
	"time"
)

type ApplicationStatus string

const (
	StatusPending  ApplicationStatus = "PENDING"
	StatusAccepted ApplicationStatus = "ACCEPTED"
	StatusRejected ApplicationStatus = "REJECTED"
)

type Application struct {
	ID                    int64             `json:"applicationId"`
	FullName              string            `json:"fullName"`
	Email                 string            `json:"email"`
	Phone                 *string           `json:"phone,omitempty"`
	OrganisationName      string            `json:"organisationName"`
	OrganisationURL       *string           `json:"organisationUrl,omitempty"`
	ProjectName           string            `json:"projectName"`
	TypeID                int64             `json:"typeId"`
	ExpectedResults       string            `json:"expectedResults"`
	IsPayed               bool              `json:"isPayed"`
	AdditionalInformation *string           `json:"additionalInformation,omitempty"`
	Status                ApplicationStatus `json:"status"`
	RejectedReason        *string           `json:"rejectedReason,omitempty"`
	ReviewerID            *int64            `json:"-"`
	CreatedAt             time.Time         `json:"-"`
	UpdatedAt             time.Time         `json:"-"`
}

type ApplicationPreview struct {
	ID               int64             `json:"externalApplicationId"`
	ProjectName      string            `json:"projectName"`
	TypeName         string            `json:"typeName"`
	Initiator        string            `json:"initiator"`
	OrganisationName string            `json:"organisationName"`
	DateUpdated      time.Time         `json:"dateUpdated"`
	Status           ApplicationStatus `json:"status"`
	RejectionMessage *string           `json:"rejectionMessage,omitempty"`
}

func (s ApplicationStatus) Valid() bool {
	return s == StatusPending || s == StatusAccepted || s == StatusRejected
}

func (a *Application) IsPending() bool {
	return a.Status == StatusPending
}

func (a *Application) Accept(reviewerID int64) error {
	if !a.IsPending() {
		return ErrApplicationNotPending
	}
	a.Status = StatusAccepted
	a.ReviewerID = &reviewerID
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Application) Reject(reviewerID int64, reason string) error {
	if !a.IsPending() {
		return ErrApplicationNotPending
	}
	a.Status = StatusRejected
	a.RejectedReason = &reason
	a.ReviewerID = &reviewerID
	a.UpdatedAt = time.Now()
	return nil
}
