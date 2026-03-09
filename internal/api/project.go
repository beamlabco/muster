package api

import (
	"fmt"
	"time"
)

// ProjectResponse represents a project in API response
type ProjectResponse struct {
	ID                int       `json:"id"`
	OrganizationID    int       `json:"organizationId"`
	Name              string    `json:"name"`
	DiscordWebhookURL *string   `json:"discordWebhookUrl"`
	SummaryEnabled    bool      `json:"summaryEnabled"`
	SummaryTime       *string   `json:"summaryTime"`
	Timezone          *string   `json:"timezone"`
	LastSummaryAt     *string   `json:"lastSummaryAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// CreateProjectRequest represents the create project request
type CreateProjectRequest struct {
	Name              string  `json:"name"`
	DiscordWebhookURL *string `json:"discordWebhookUrl,omitempty"`
	SummaryEnabled    *bool   `json:"summaryEnabled,omitempty"`
	SummaryTime       *string `json:"summaryTime,omitempty"`
	Timezone          *string `json:"timezone,omitempty"`
}

// UpdateProjectRequest represents the update project request
type UpdateProjectRequest struct {
	Name              *string `json:"name,omitempty"`
	DiscordWebhookURL *string `json:"discordWebhookUrl,omitempty"`
	SummaryEnabled    *bool   `json:"summaryEnabled,omitempty"`
	SummaryTime       *string `json:"summaryTime,omitempty"`
	Timezone          *string `json:"timezone,omitempty"`
}

// AddMembersRequest represents the add members request
type AddMembersRequest struct {
	UserIDs []int `json:"userIds"`
}

// ProjectMemberResponse represents a project member
type ProjectMemberResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetMyProjects retrieves the current user's projects
func (c *Client) GetMyProjects() ([]*ProjectResponse, error) {
	var result []*ProjectResponse
	resp, err := c.http.R().
		SetResult(&result).
		Get("/api/projects/my")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return result, nil
}

// GetProjects retrieves all projects in the organization
func (c *Client) GetProjects() ([]*ProjectResponse, error) {
	var result []*ProjectResponse
	resp, err := c.http.R().
		SetResult(&result).
		Get("/api/projects")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return result, nil
}

// GetProject retrieves a single project by ID
func (c *Client) GetProject(projectID int) (*ProjectResponse, error) {
	var result ProjectResponse
	resp, err := c.http.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/projects/%d", projectID))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return &result, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(req *CreateProjectRequest) (*ProjectResponse, error) {
	var result ProjectResponse
	resp, err := c.http.R().
		SetBody(req).
		SetResult(&result).
		Post("/api/projects")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return &result, nil
}

// UpdateProject updates project settings
func (c *Client) UpdateProject(projectID int, req *UpdateProjectRequest) (*ProjectResponse, error) {
	var result ProjectResponse
	resp, err := c.http.R().
		SetBody(req).
		SetResult(&result).
		Patch(fmt.Sprintf("/api/projects/%d", projectID))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return &result, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID int) error {
	resp, err := c.http.R().
		Delete(fmt.Sprintf("/api/projects/%d", projectID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp)
	}

	return nil
}

// AddProjectMembers adds members to a project
func (c *Client) AddProjectMembers(projectID int, userIDs []int) error {
	resp, err := c.http.R().
		SetBody(&AddMembersRequest{UserIDs: userIDs}).
		Post(fmt.Sprintf("/api/projects/%d/members", projectID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp)
	}

	return nil
}

// RemoveProjectMember removes a member from a project
func (c *Client) RemoveProjectMember(projectID, userID int) error {
	resp, err := c.http.R().
		Delete(fmt.Sprintf("/api/projects/%d/members/%d", projectID, userID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp)
	}

	return nil
}

// GetProjectMembers retrieves project members
func (c *Client) GetProjectMembers(projectID int) ([]*ProjectMemberResponse, error) {
	var result []*ProjectMemberResponse
	resp, err := c.http.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/projects/%d/members", projectID))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp)
	}

	return result, nil
}
