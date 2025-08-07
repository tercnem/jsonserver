package config

type APIConfig struct {
	Port any `json:"port"`
}

// "method": "POST",
//
//	"status": 200,
//	"path": "/api/auth/isTokenValid",
//	"jsonPath": "./isvalid.json"

var ApiConfig *APIConfig
