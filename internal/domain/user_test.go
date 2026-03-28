package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{name: "admin", role: RoleAdmin, expected: true},
		{name: "user", role: RoleUser, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.expected, user.IsAdmin())
		})
	}
}

func TestRole_Valid(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{name: "admin", role: RoleAdmin, expected: true},
		{name: "user", role: RoleUser, expected: true},
		{name: "invalid", role: Role("INVALID"), expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.Valid())
		})
	}
}
