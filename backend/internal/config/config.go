package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config representa a configuração global da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	SSO      SSOConfig
	Env      string
}

// ServerConfig contém configurações do servidor HTTP
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig contém configurações do banco de dados
type DatabaseConfig struct {
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string
	MongoURI         string
}

// AuthConfig contém configurações de autenticação
type AuthConfig struct {
	JwtSecret          string
	JwtExpiryHours     int
	RefreshSecret      string
	RefreshExpiryHours int
}

// SSOConfig contém as configurações para o cliente OAuth2 SSO
type SSOConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

// Load carrega configurações a partir de variáveis de ambiente
func Load() (*Config, error) {
	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("porta do servidor inválida: %w", err)
	}

	pgPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("porta do postgres inválida: %w", err)
	}

	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("JWT expiry inválido: %w", err)
	}

	refreshExpiryHours, err := strconv.Atoi(getEnv("REFRESH_EXPIRY_HOURS", "168"))
	if err != nil {
		return nil, fmt.Errorf("Refresh token expiry inválido: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port:         serverPort,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
			PostgresPort:     pgPort,
			PostgresUser:     getEnv("POSTGRES_USER", "cecor"),
			PostgresPassword: getEnv("POSTGRES_PASSWORD", "cecor123"),
			PostgresDB:       getEnv("POSTGRES_DB", "cepr_db"),
			PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
			MongoURI:         getEnv("MONGO_URI", "mongodb://mongo:27017"),
		},
		Auth: AuthConfig{
			JwtSecret:          getEnv("JWT_SECRET", "sua_chave_secreta_muito_segura"), // WARNING: Default value for development only. Do not use in production.
			JwtExpiryHours:     jwtExpiryHours,
			RefreshSecret:      getEnv("REFRESH_SECRET", "outra_chave_secreta_muito_segura"), // WARNING: Default value for development only. Do not use in production.
			RefreshExpiryHours: refreshExpiryHours,
		},
		SSO: SSOConfig{
			ClientID:     getEnv("SSO_CLIENT_ID", "cecor-frontend"),
			ClientSecret: getEnv("SSO_CLIENT_SECRET", "cecor-secret"), // WARNING: Default value for development only. Do not use in production.
			RedirectURL:  getEnv("SSO_REDIRECT_URL", "http://localhost:4200/*"),
			AuthURL:      getEnv("SSO_AUTH_URL", "http://localhost:8081/realms/cecor/protocol/openid-connect/auth"),
			TokenURL:     getEnv("SSO_TOKEN_URL", "http://localhost:8081/realms/cecor/protocol/openid-connect/token"),
			UserInfoURL:  getEnv("SSO_USER_INFO_URL", "http://localhost:8081/realms/cecor/protocol/openid-connect/userinfo"),
		},
		Env: getEnv("APP_ENV", "development"),
	}, nil
}

// getEnv obtém valor de variável de ambiente ou retorna valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
