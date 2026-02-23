// backend/internal/api/handlers/auth_handler.go
package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/devdavidalonso/cecor/backend/internal/api/middleware"
	"github.com/devdavidalonso/cecor/backend/internal/config"
	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/service/users"
	"github.com/devdavidalonso/cecor/backend/pkg/errors"
	"golang.org/x/oauth2"
)

const ssoStateCookie = "sso_state"

func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}

// AuthResponse representa a resposta de uma autenticação bem-sucedida
type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken"`
	User         models.User `json:"user"`
}

// UserInfo represents the user information returned from the SSO's userinfo endpoint
type UserInfo struct {
	UserID      string `json:"sub"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

// AuthHandler implementa os handlers HTTP para autenticação
type AuthHandler struct {
	userService users.Service
	cfg         *config.Config
	ssoConfig   *oauth2.Config
}

// NewAuthHandler cria uma nova instância de AuthHandler
func NewAuthHandler(userService users.Service, cfg *config.Config, ssoConfig *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		cfg:         cfg,
		ssoConfig:   ssoConfig,
	}
}

// SSOLogin redirects the user to the SSO provider for authentication.
func (h *AuthHandler) SSOLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		errors.RespondWithError(w, http.StatusInternalServerError, "Failed to generate state for SSO")
		return
	}

	// Store the state in a short-lived cookie
	http.SetCookie(w, &http.Cookie{
		Name:     ssoStateCookie,
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
	})

	// Redirect user to consent page to ask for permission
	url := h.ssoConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// SSOCallback handles the callback from the SSO provider.
func (h *AuthHandler) SSOCallback(w http.ResponseWriter, r *http.Request) {
	// Check state cookie
	stateCookie, err := r.Cookie(ssoStateCookie)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "SSO state cookie not found")
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid SSO state")
		return
	}

	// Exchange authorization code for a token
	code := r.URL.Query().Get("code")
	fmt.Printf("Attempting to exchange code: %s\n", code)
	token, err := h.ssoConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Error exchanging token: %v\n", err)
		errors.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange token with SSO provider")
		return
	}
	fmt.Printf("Token exchanged successfully. Access Token: %s...\n", token.AccessToken[:10])

	// Use the token to get user info
	userInfo, err := h.getUserInfo(token)
	if err != nil {
		errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Extract the best matching role
	var role string
	if len(userInfo.RealmAccess.Roles) > 0 {
		// Priority list: admin/manager > professor > student/guardian
		for _, r := range userInfo.RealmAccess.Roles {
			normalized := normalizeRole(r)
			switch normalized {
			case "administrator", "admin", "gestor", "teacher", "professor", "student", "aluno", "responsável", "responsavel":
				role = normalized
			}
			if role != "" {
				break
			}
		}
	}
	if role == "" {
		role = "user"
	}

	// At this point, you have the user's email from the SSO provider.
	// You can now find or create a user in your local database.
	user, err := h.userService.FindOrCreateByEmail(context.Background(), userInfo.Email, userInfo.Name, role)
	if err != nil {
		errors.RespondWithError(w, http.StatusInternalServerError, "Failed to process user information")
		return
	}

	// Keycloak é a única fonte de autenticação/token.
	// Não geramos JWT local aqui; apenas provisionamos o usuário local e redirecionamos
	// para o frontend iniciar o fluxo OIDC direto.
	_ = user
	redirectURL := "http://localhost:4201/auth/login"
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) getUserInfo(token *oauth2.Token) (*UserInfo, error) {
	client := h.ssoConfig.Client(context.Background(), token)
	resp, err := client.Get(h.cfg.SSO.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from SSO provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SSO provider returned non-200 status for user info: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response body: %w", err)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Verify retorna as informações do usuário autenticado
func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		errors.RespondWithError(w, http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"valid": true,
		"user":  claims,
	})
}
