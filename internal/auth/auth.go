package auth

import (
	"fmt"
	"strings"

	"github.com/muster/cli/internal/api"
	"github.com/muster/cli/internal/config"
)

// Service handles authentication operations
type Service struct {
	client *api.Client
	config *config.Config
}

// NewService creates a new authentication service
func NewService(client *api.Client, cfg *config.Config) *Service {
	return &Service{
		client: client,
		config: cfg,
	}
}

// Register creates a new user and organization
func (s *Service) Register(email, password, name, organizationName string) error {
	// Validate inputs
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	if err := validateName(name); err != nil {
		return err
	}
	if err := validateOrganizationName(organizationName); err != nil {
		return err
	}

	// Call API
	req := &api.RegisterRequest{
		Email:            email,
		Password:         password,
		Name:             name,
		OrganizationName: organizationName,
	}

	resp, err := s.client.Register(req)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	// Save to config
	s.config.SetAuthData(
		resp.Token,
		resp.User.ToConfigUser(),
		resp.Organization.ToConfigOrganization(),
	)

	if err := config.Save(s.config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Update client token
	s.client.SetToken(resp.Token)

	return nil
}

// Join registers a user via an invitation token
func (s *Service) Join(email, password, name, token string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	if err := validateName(name); err != nil {
		return err
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("invitation token is required")
	}

	req := &api.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
		Token:    token,
	}

	resp, err := s.client.Register(req)
	if err != nil {
		return fmt.Errorf("join failed: %w", err)
	}

	s.config.SetAuthData(
		resp.Token,
		resp.User.ToConfigUser(),
		resp.Organization.ToConfigOrganization(),
	)

	if err := config.Save(s.config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	s.client.SetToken(resp.Token)

	return nil
}

// Login authenticates a user
func (s *Service) Login(email, password string) error {
	// Validate inputs
	if err := validateEmail(email); err != nil {
		return err
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}

	// Call API
	req := &api.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := s.client.Login(req)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Save to config
	s.config.SetAuthData(
		resp.Token,
		resp.User.ToConfigUser(),
		resp.Organization.ToConfigOrganization(),
	)

	if err := config.Save(s.config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Update client token
	s.client.SetToken(resp.Token)

	return nil
}

// Logout clears the authentication data
func (s *Service) Logout() error {
	s.config.Clear()
	if err := config.Save(s.config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	s.client.SetToken("")
	return nil
}

// IsAuthenticated checks if the user is authenticated
func (s *Service) IsAuthenticated() bool {
	return s.config.IsAuthenticated()
}

// GetUser returns the authenticated user
func (s *Service) GetUser() *config.User {
	return s.config.User
}

// GetOrganization returns the user's organization
func (s *Service) GetOrganization() *config.Organization {
	return s.config.Organization
}

// Validation functions

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

func validateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(name) > 255 {
		return fmt.Errorf("name must be less than 255 characters")
	}
	return nil
}

func validateOrganizationName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return fmt.Errorf("organization name must be at least 2 characters")
	}
	if len(name) > 255 {
		return fmt.Errorf("organization name must be less than 255 characters")
	}
	return nil
}
