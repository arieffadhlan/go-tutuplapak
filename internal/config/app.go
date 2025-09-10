package config

import (
	"fmt"
	"os"
)

type App struct {
	Port      string
	Envs      string
	JWTSecret string
}

func loadAppConfig() App {
	return App{
		Port:      getAppPort(),
		Envs:      getAppEnvs(),
		JWTSecret: getAppJWTSecret(),
	}
}

func getAppPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	return fmt.Sprintf(":%s", port)
}

func getAppEnvs() string {
	env := os.Getenv("APP_ENVS")
	if env == "" {
		return "development"
	}
	return env
}

func getAppJWTSecret() string {
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		return "default-secret"
	}
	return jwt
}
