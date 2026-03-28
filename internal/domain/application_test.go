package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplication_Accept(t *testing.T) {
	reviewerID := int64(1)

	tests := []struct {
		name    string
		status  ApplicationStatus
		wantErr bool
		errType error
	}{
		{
			name:    "accept pending application",
			status:  StatusPending,
			wantErr: false,
		},
		{
			name:    "accept already accepted application",
			status:  StatusAccepted,
			wantErr: true,
			errType: ErrApplicationNotPending,
		},
		{
			name:    "accept rejected application",
			status:  StatusRejected,
			wantErr: true,
			errType: ErrApplicationNotPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				ID:        1,
				Status:    tt.status,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := app.Accept(reviewerID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, StatusAccepted, app.Status)
				assert.Equal(t, &reviewerID, app.ReviewerID)
			}
		})
	}
}

func TestApplication_Reject(t *testing.T) {
	reviewerID := int64(1)
	reason := "test reason"

	tests := []struct {
		name    string
		status  ApplicationStatus
		reason  string
		wantErr bool
		errType error
	}{
		{
			name:    "reject pending application",
			status:  StatusPending,
			reason:  reason,
			wantErr: false,
		},
		{
			name:    "reject already accepted application",
			status:  StatusAccepted,
			reason:  reason,
			wantErr: true,
			errType: ErrApplicationNotPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				ID:        1,
				Status:    tt.status,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := app.Reject(reviewerID, tt.reason)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, StatusRejected, app.Status)
				assert.Equal(t, &reviewerID, app.ReviewerID)
				assert.Equal(t, &reason, app.RejectedReason)
			}
		})
	}
}

func TestApplication_IsPending(t *testing.T) {
	tests := []struct {
		name     string
		status   ApplicationStatus
		expected bool
	}{
		{name: "pending", status: StatusPending, expected: true},
		{name: "accepted", status: StatusAccepted, expected: false},
		{name: "rejected", status: StatusRejected, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{Status: tt.status}
			assert.Equal(t, tt.expected, app.IsPending())
		})
	}
}
