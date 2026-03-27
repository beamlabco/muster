package user

import (
	"fmt"
	"strings"

	"github.com/muster/cli/internal/api"
)

// Valid assignable roles (primary cannot be assigned)
var ValidRoles = []string{"manager", "member"}

// Service handles user operations
type Service struct {
	client *api.Client
}

// NewService creates a new user service
func NewService(client *api.Client) *Service {
	return &Service{
		client: client,
	}
}

// GetMembers retrieves all members of the organization
func (s *Service) GetMembers() ([]*api.MemberResponse, error) {
	resp, err := s.client.GetMembers()
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	return resp.Users, nil
}

// UpdateRole updates a user's role
func (s *Service) UpdateRole(userID int, role string) (*api.UserResponse, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if !isValidRole(role) {
		return nil, fmt.Errorf("invalid role, must be one of: %s", strings.Join(ValidRoles, ", "))
	}

	resp, err := s.client.UpdateUserRole(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	return resp, nil
}

// ResetPassword resets another user's password (primary only)
func (s *Service) ResetPassword(userID int, newPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters")
	}

	err := s.client.ResetUserPassword(userID, newPassword)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}

// ChangePassword changes the authenticated user's password
func (s *Service) ChangePassword(currentPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters")
	}

	err := s.client.ChangePassword(currentPassword, newPassword)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

func isValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}
