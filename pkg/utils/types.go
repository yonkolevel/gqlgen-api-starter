package utils

import (
	"strings"
)

// ContextKey defines a type for context keys shared in the app
type ContextKey string

// ServerConfig defines the configuration for the server
type ServerConfig struct {
	Host          string
	Port          string
	URISchema     string
	Version       string
	SessionSecret string
	JWT           JWTConfig
	GraphQL       GQLConfig
	Database      DBConfig
	AuthProviders []AuthProvider
	Spaces        SpacesConfig
}

type SpacesConfig struct {
	Key      string
	Secret   string
	Endpoint string
}

//JWTConfig defines the options for JWT tokens
type JWTConfig struct {
	Secret    string
	Algorithm string
}

// GQLConfig defines the configuration for the GQL Server
type GQLConfig struct {
	Path                string
	PlaygroundPath      string
	IsPlaygroundEnabled bool
}

// DBConfig defines the configuration for the DB config
type DBConfig struct {
	Dialect     string
	DSN         string
	SeedDB      bool
	LogMode     bool
	AutoMigrate bool
}

// AuthProvider defines the configuration for the Goth config
type AuthProvider struct {
	Provider  string
	ClientKey string
	Secret    string
	Domain    string // If needed, like with auth0
	Scopes    []string
}

// ListenEndpoint builds the endpoint string (host + port)
func (s *ServerConfig) ListenEndpoint() string {
	if s.Port == "80" {
		return s.Host
	}
	return s.Host + ":" + s.Port
}

// VersionedEndpoint builds the endpoint string (host + port + version)
func (s *ServerConfig) VersionedEndpoint(path string) string {
	return "/" + s.Version + path
}

// SchemaVersionedEndpoint builds the schema endpoint string (schema + host + port + version)
func (s *ServerConfig) SchemaVersionedEndpoint(path string) string {
	if s.Port == "80" {
		return s.URISchema + s.Host + "/" + s.Version + path
	}
	return s.URISchema + s.Host + ":" + s.Port + "/" + s.Version + path
}

func NewServerConfig() *ServerConfig {
	var serverconf = &ServerConfig{
		Host:          MustGet("SERVER_HOST"),
		Port:          MustGet("PORT"),
		URISchema:     MustGet("SERVER_URI_SCHEMA"),
		Version:       MustGet("SERVER_PATH_VERSION"),
		SessionSecret: MustGet("SESSION_SECRET"),
		JWT: JWTConfig{
			Secret:    MustGet("AUTH_JWT_SECRET"),
			Algorithm: MustGet("AUTH_JWT_SIGNING_ALGORITHM"),
		},
		GraphQL: GQLConfig{
			Path:                MustGet("GQL_SERVER_GRAPHQL_PATH"),
			PlaygroundPath:      MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH"),
			IsPlaygroundEnabled: MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED"),
		},
		Database: DBConfig{
			Dialect:     MustGet("GORM_DIALECT"),
			DSN:         MustGet("DATABASE_URL"),
			SeedDB:      MustGetBool("GORM_SEED_DB"),
			LogMode:     MustGetBool("GORM_LOGMODE"),
			AutoMigrate: MustGetBool("GORM_AUTOMIGRATE"),
		},
		AuthProviders: []AuthProvider{
			{
				Provider:  "google",
				ClientKey: MustGet("PROVIDER_GOOGLE_KEY"),
				Secret:    MustGet("PROVIDER_GOOGLE_SECRET"),
			},
			{
				Provider:  "auth0",
				ClientKey: MustGet("PROVIDER_AUTH0_KEY"),
				Secret:    MustGet("PROVIDER_AUTH0_SECRET"),
				Domain:    MustGet("PROVIDER_AUTH0_DOMAIN"),
				Scopes:    strings.Split(MustGet("PROVIDER_AUTH0_SCOPES"), ","),
			},
			{
				Provider:  "facebook",
				ClientKey: MustGet("PROVIDER_FACEBOOK_KEY"),
				Secret:    MustGet("PROVIDER_FACEBOOK_SECRET"),
			},
			{
				Provider:  "twitter",
				ClientKey: MustGet("PROVIDER_TWITTER_KEY"),
				Secret:    MustGet("PROVIDER_TWITTER_SECRET"),
			},
		},
		Spaces: SpacesConfig{
			Key:      MustGet("SPACES_KEY"),
			Secret:   MustGet("SPACES_SECRET"),
			Endpoint: MustGet("SPACES_ENDPOINT"),
		},
	}

	return serverconf
}

var TestServerconf = &ServerConfig{
	Host:          "localhost",
	Port:          "7777",
	URISchema:     "http",
	Version:       "1",
	SessionSecret: "secret",
	JWT: JWTConfig{
		Secret:    "secret",
		Algorithm: "HS512",
	},
	GraphQL: GQLConfig{
		Path:                "/graphql",
		PlaygroundPath:      "/playground",
		IsPlaygroundEnabled: false,
	},
	Database: DBConfig{
		Dialect:     "postgres",
		DSN:         "postgresql://db:1233@testdomain.com:25060/dev-db",
		SeedDB:      false,
		LogMode:     true,
		AutoMigrate: false,
	},
	AuthProviders: []AuthProvider{
		{
			Provider:  "test-auth-provider",
			ClientKey: "key",
			Secret:    "secret",
		},
	},
	Spaces: SpacesConfig{
		Key:      "key",
		Secret:   "secret",
		Endpoint: "",
	},
}
