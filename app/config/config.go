package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	JwtKey string
	AppName string
	AppEnv string
	AppPort uint
	SentryDSN string
	Subjects string
	RefreshExpiry time.Duration
	TokenExpiry time.Duration
	DefaultPerPage uint64
	// Services
	UserServiceHost string
	UserServicePort uint
	SessionServiceHost string
	SessionServicePort uint
	ResumeServiceHost string
	ResumeServicePort uint
}
const (
	UserKey = "USER_ID"
	UserSessionKey = "USER_SESSION_ID"
)

func Configure() *Config {
	usrPort, _ := strconv.Atoi(os.Getenv("USER_SERVICE_PORT"))
	sessPort, _ := strconv.Atoi(os.Getenv("SESSION_SERVICE_PORT"))
	appPort, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	resumePort, _ := strconv.Atoi(os.Getenv("RESUME_SERVICE_PORT"))

	return &Config{
		JwtKey: os.Getenv("JWT_KEY"),
		AppName: os.Getenv("SERVICE_NAME"),
		AppEnv: os.Getenv("APP_ENV"),
		AppPort: uint(appPort),
		SentryDSN: os.Getenv("SENTRY_DSN"),
		Subjects: "me,resume,review",
		TokenExpiry: time.Hour,
		RefreshExpiry: time.Hour * 24,

		UserServiceHost: os.Getenv("USER_SERVICE_HOST"),
		UserServicePort: uint(usrPort),
		DefaultPerPage: 50,

		SessionServiceHost: os.Getenv("SESSION_SERVICE_HOST"),
		SessionServicePort: uint(sessPort),
		ResumeServiceHost: os.Getenv("RESUME_SERVICE_HOST"),
		ResumeServicePort: uint(resumePort),

	}
}