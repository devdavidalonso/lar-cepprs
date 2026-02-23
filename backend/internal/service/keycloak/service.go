package keycloak

import (
	"context"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

// KeycloakService handles Keycloak Admin API operations
type KeycloakService struct {
	client      *gocloak.GoCloak
	realm       string
	adminRealm  string
	clientID    string
	username    string
	password    string
	accessToken string
}

// NewKeycloakService creates a new Keycloak service instance
func NewKeycloakService() *KeycloakService {
	baseURL := os.Getenv("KEYCLOAK_ADMIN_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	return &KeycloakService{
		client:     gocloak.NewClient(baseURL),
		realm:      os.Getenv("KEYCLOAK_TARGET_REALM"),
		adminRealm: os.Getenv("KEYCLOAK_ADMIN_REALM"),
		clientID:   os.Getenv("KEYCLOAK_ADMIN_CLIENT_ID"),
		username:   os.Getenv("KEYCLOAK_ADMIN_USERNAME"),
		password:   os.Getenv("KEYCLOAK_ADMIN_PASSWORD"),
	}
}

// authenticate obtains an admin access token
func (s *KeycloakService) authenticate(ctx context.Context) error {
	token, err := s.client.LoginAdmin(ctx, s.username, s.password, s.adminRealm)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Keycloak: %w", err)
	}
	s.accessToken = token.AccessToken
	return nil
}

// CreateUserRequest represents the data needed to create a user in Keycloak
type CreateUserRequest struct {
	Username      string
	Email         string
	FirstName     string
	LastName      string
	Enabled       bool
	EmailVerified bool
}

// CreateUser creates a new user in Keycloak and returns the user ID
func (s *KeycloakService) CreateUser(ctx context.Context, req CreateUserRequest) (string, error) {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return "", err
	}

	// Prepare user object
	enabled := req.Enabled
	emailVerified := req.EmailVerified
	user := gocloak.User{
		Username:        gocloak.StringP(req.Username),
		Email:           gocloak.StringP(req.Email),
		FirstName:       gocloak.StringP(req.FirstName),
		LastName:        gocloak.StringP(req.LastName),
		Enabled:         &enabled,
		EmailVerified:   &emailVerified,
		RequiredActions: &[]string{},
	}

	// Create user in Keycloak
	userID, err := s.client.CreateUser(ctx, s.accessToken, s.realm, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user in Keycloak: %w", err)
	}

	return userID, nil
}

// AssignRole assigns a realm role to a user
func (s *KeycloakService) AssignRole(ctx context.Context, userID, roleName string) error {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return err
	}

	// Get the role by name
	role, err := s.client.GetRealmRole(ctx, s.accessToken, s.realm, roleName)
	if err != nil {
		return fmt.Errorf("failed to get role '%s': %w", roleName, err)
	}

	// Assign role to user
	err = s.client.AddRealmRoleToUser(ctx, s.accessToken, s.realm, userID, []gocloak.Role{*role})
	if err != nil {
		return fmt.Errorf("failed to assign role '%s' to user: %w", roleName, err)
	}

	return nil
}

// DeleteUser deletes a user from Keycloak
func (s *KeycloakService) DeleteUser(ctx context.Context, userID string) error {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return err
	}

	err := s.client.DeleteUser(ctx, s.accessToken, s.realm, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user from Keycloak: %w", err)
	}

	return nil
}

// DisableUser disables a user in Keycloak
func (s *KeycloakService) DisableUser(ctx context.Context, userID string) error {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return err
	}

	// Get user
	user, err := s.client.GetUserByID(ctx, s.accessToken, s.realm, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Disable user
	enabled := false
	user.Enabled = &enabled

	err = s.client.UpdateUser(ctx, s.accessToken, s.realm, *user)
	if err != nil {
		return fmt.Errorf("failed to disable user: %w", err)
	}

	return nil
}

// EnableUser enables a user in Keycloak
func (s *KeycloakService) EnableUser(ctx context.Context, userID string) error {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return err
	}

	// Get user
	user, err := s.client.GetUserByID(ctx, s.accessToken, s.realm, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Enable user
	enabled := true
	user.Enabled = &enabled

	err = s.client.UpdateUser(ctx, s.accessToken, s.realm, *user)
	if err != nil {
		return fmt.Errorf("failed to enable user: %w", err)
	}

	return nil
}

// SetTemporaryPassword sets a temporary password for a user
func (s *KeycloakService) SetTemporaryPassword(ctx context.Context, userID, password string) error {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return err
	}

	temporary := false
	err := s.client.SetPassword(ctx, s.accessToken, userID, s.realm, password, temporary)
	if err != nil {
		return fmt.Errorf("failed to set temporary password: %w", err)
	}

	return nil
}

// GetUsersByRole fetches users who have a specific realm role
func (s *KeycloakService) GetUsersByRole(ctx context.Context, roleName string) ([]*gocloak.User, error) {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return nil, err
	}

	// Get users with the specified role
	users, err := s.client.GetUsersByRoleName(ctx, s.accessToken, s.realm, roleName, gocloak.GetUsersByRoleParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users with role '%s': %w", roleName, err)
	}

	return users, nil
}

// AddUserToGroup adds a user to a Keycloak group
func (s *KeycloakService) AddUserToGroup(ctx context.Context, userID, groupName string) error {
	// Authenticate first
	if err := s.authenticate(ctx); err != nil {
		return err
	}

	// Get groups to find the ID of the group with the given name
	groups, err := s.client.GetGroups(ctx, s.accessToken, s.realm, gocloak.GetGroupsParams{
		Search: gocloak.StringP(groupName),
	})
	if err != nil {
		return fmt.Errorf("failed to search for group '%s': %w", groupName, err)
	}

	var groupID string
	for _, g := range groups {
		if g.Name != nil && *g.Name == groupName {
			groupID = *g.ID
			break
		}
	}

	if groupID == "" {
		return fmt.Errorf("group '%s' not found", groupName)
	}

	// Add user to group
	err = s.client.AddUserToGroup(ctx, s.accessToken, s.realm, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to add user to group '%s': %w", groupName, err)
	}

	return nil
}
