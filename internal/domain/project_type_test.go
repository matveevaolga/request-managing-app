package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectType_Validate(t *testing.T) {
	tests := []struct {
		name        string
		projectType *ProjectType
		wantErr     bool
		errType     error
	}{
		{
			name:        "valid project type name",
			projectType: &ProjectType{Name: "Startup"},
			wantErr:     false,
		},
		{
			name:        "empty name",
			projectType: &ProjectType{Name: ""},
			wantErr:     true,
			errType:     ErrInvalidProjectTypeName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.projectType.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
