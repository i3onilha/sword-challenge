package models

import "testing"

func TestUser_IsManager(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "User is manager",
			role:     RoleManager,
			expected: true,
		},
		{
			name:     "User is technician",
			role:     RoleTechnician,
			expected: false,
		},
		{
			name:     "User has unknown role",
			role:     UserRole("other"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Role: tt.role}
			if got := u.IsManager(); got != tt.expected {
				t.Errorf("IsManager() = %v, want %v", got, tt.expected)
			}
		})
	}
}
func TestUser_IsTechnician(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "User is technician",
			role:     RoleTechnician,
			expected: true,
		},
		{
			name:     "User is manager",
			role:     RoleManager,
			expected: false,
		},
		{
			name:     "User has unknown role",
			role:     UserRole("other"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Role: tt.role}
			if got := u.IsTechnician(); got != tt.expected {
				t.Errorf("IsTechnician() = %v, want %v", got, tt.expected)
			}
		})
	}
}
