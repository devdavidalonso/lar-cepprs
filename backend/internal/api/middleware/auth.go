package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/devdavidalonso/cecor/backend/internal/auth"
	"github.com/devdavidalonso/cecor/backend/internal/config"
	"github.com/devdavidalonso/cecor/backend/pkg/errors"
)

// UserClaims armazena informações do usuário no contexto
type UserClaims struct {
	UserID int64    `json:"userId"`
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

// contextKey é um tipo para chaves de contexto
type contextKey string

// userClaimsKey é a chave para armazenar claims do usuário no contexto
const userClaimsKey contextKey = "userClaims"

// Authenticate verifica o token JWT de autenticação
func Authenticate(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrair token do cabeçalho Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errors.RespondWithError(w, http.StatusUnauthorized, "Token de autenticação não fornecido")
				return
			}

			// Verificar formato do token (Bearer token)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				errors.RespondWithError(w, http.StatusUnauthorized, "Formato de token inválido")
				return
			}

			// Validar token (OIDC Keycloak)
			// SIMULATION MODE: Em desenvolvimento, aceitar tokens de simulação
			if cfg.Env == "development" && strings.HasPrefix(parts[1], "simulation-") {
				// Simular usuário baseado no token
				var userClaims *UserClaims
				if parts[1] == "simulation-admin" {
					userClaims = &UserClaims{
						Name:  "Admin Simulado",
						Email: "admin@cecor.org",
						Roles: []string{"admin"},
					}
				} else if parts[1] == "simulation-professor" {
					userClaims = &UserClaims{
						Name:  "Professor Simulado",
						Email: "professor@cecor.org",
						Roles: []string{"professor"},
					}
				} else {
					errors.RespondWithError(w, http.StatusUnauthorized, "Token de simulação desconhecido")
					return
				}

				fmt.Printf("⚠️  Modo Simulação: Autenticado como %s (%s)\n", userClaims.Name, userClaims.Roles[0])
				ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			claims, err := auth.ValidateOIDCToken(r.Context(), parts[1])
			if err != nil {
				fmt.Printf("❌ Authentication Failed: %v\n", err)
				errors.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Token inválido: %v", err))
				return
			}

			// Extrair claims específicas do usuário
			userClaims := &UserClaims{
				Email: claims["email"].(string),
				Name:  claims["name"].(string),
			}

			// Tentar extrair UserID se disponível (pode não estar no token do Keycloak inicialmente)
			// No Keycloak, o ID do usuário é 'sub' (UUID), mas nosso sistema usa int64.
			// Por enquanto, vamos deixar UserID zerado ou tentar mapear via email se necessário.
			// Para o MVP, vamos confiar no email como identificador principal se o userId não for compatível.
			if sub, ok := claims["sub"].(string); ok {
				// TODO: Mapear UUID do Keycloak para ID local ou usar UUID no sistema todo
				// Por enquanto, não vamos quebrar se não conseguir converter
				fmt.Printf("User authenticated: %s (%s)\n", userClaims.Email, sub)
			}

			// Extrair roles do realm_access
			if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
				if roles, ok := realmAccess["roles"].([]interface{}); ok {
					for _, role := range roles {
						if roleStr, ok := role.(string); ok {
							userClaims.Roles = append(userClaims.Roles, roleStr)
						}
					}
				}
			}

			// Adicionar claims ao contexto
			ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)

			// Chamar o próximo handler com o contexto atualizado
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin verifica se o usuário tem papel de administrador
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter claims do contexto
		claims, ok := r.Context().Value(userClaimsKey).(*UserClaims)
		if !ok {
			errors.RespondWithError(w, http.StatusUnauthorized, "Usuário não autenticado")
			return
		}

		// Verificar se o usuário tem papel de administrador
		isAdmin := false
		for _, role := range claims.Roles {
			normalizedRole := strings.ToLower(strings.TrimSpace(role))
			if normalizedRole == "admin" || normalizedRole == "administrator" || normalizedRole == "administrador" || normalizedRole == "gestor" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			errors.RespondWithError(w, http.StatusForbidden, "Acesso negado: privilégios de administrador necessários")
			return
		}

		// Usuário tem permissão, prosseguir
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extrai informações do usuário do contexto
func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*UserClaims)
	return claims, ok
}
