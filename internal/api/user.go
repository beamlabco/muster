package api

import (
	"fmt"
)

// UpdateRoleRequest represents the role update request
type UpdateRoleRequest struct {
	Role string `json:"role"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// GetMembersResponse represents the get members response
type GetMembersResponse struct {
	Users []*MemberResponse `json:"users"`
}

// MemberResponse represents a team member in API response
type MemberResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetMembers retrieves all members of the current organization
func (c *Client) GetMembers() (*GetMembersResponse, error) {
	var result GetMembersResponse
	resp, err := c.http.R().
		SetResult(&result).
		Get("/api/users")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return &result, nil
}

// ChangePassword changes the authenticated user's password
func (c *Client) ChangePassword(currentPassword, newPassword string) error {
	resp, err := c.http.R().
		SetBody(&ChangePasswordRequest{
			CurrentPassword: currentPassword,
			NewPassword:     newPassword,
		}).
		Patch("/api/users/password")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp)
	}

	return nil
}

// ResetUserPasswordRequest represents the reset password request
type ResetUserPasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

// ResetUserPassword resets another user's password (primary only)
func (c *Client) ResetUserPassword(userID int, newPassword string) error {
	resp, err := c.http.R().
		SetBody(&ResetUserPasswordRequest{NewPassword: newPassword}).
		Patch(fmt.Sprintf("/api/users/%d/reset-password", userID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp)
	}

	return nil
}

// UpdateUserRole updates a user's role
func (c *Client) UpdateUserRole(userID int, role string) (*UserResponse, error) {
	var result UserResponse
	resp, err := c.http.R().
		SetBody(&UpdateRoleRequest{Role: role}).
		SetResult(&result).
		Patch(fmt.Sprintf("/api/users/%d/role", userID))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return &result, nil
}
