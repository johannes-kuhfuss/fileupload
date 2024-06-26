package config

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sanitize/sanitize"
	"github.com/johannes-kuhfuss/fileupload/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Server struct {
		Host                 string `envconfig:"SERVER_HOST"`
		Port                 string `envconfig:"SERVER_PORT" default:"8080"`
		TlsPort              string `envconfig:"SERVER_TLS_PORT" default:"8443"`
		GracefulShutdownTime int    `envconfig:"GRACEFUL_SHUTDOWN_TIME" default:"10"`
		UseTls               bool   `envconfig:"USE_TLS" default:"false"`
		CertFile             string `envconfig:"CERT_FILE" default:"./cert/cert.pem"`
		KeyFile              string `envconfig:"KEY_FILE" default:"./cert/cert.key"`
		CookieSecret         string `envconfig:"COOKIE_SECRET" default:"veryverysecret"`
	}
	Gin struct {
		Mode         string `envconfig:"GIN_MODE" default:"release"`
		TemplatePath string `envconfig:"TEMPLATE_PATH" default:"./templates/"`
	}
	Upload struct {
		Path              string            `envconfig:"UPLOAD_PATH" default:"C:\\TEMP"`
		AllowedExtensions []string          `envconfig:"ALLOWED_EXTENESION" default:".mp3,.m4a,.wav"`
		Users             map[string]string `envconfig:"USERS"`
		WriteLog          bool              `envconfig:"WRITE_LOG" default:"true"`
		LogFile           string            `envconfig:"LOG_FILE" default:"upload_log.txt"`
	}
	RunTime struct {
		Router     *gin.Engine
		ListenAddr string
		StartDate  time.Time
		Sani       *sanitize.Sanitizer
		UploadList []domain.Upload
	}
}

var (
	EnvFile = ".env"
)

func InitConfig(file string, config *AppConfig) api_error.ApiErr {
	logger.Info(fmt.Sprintf("Initalizing configuration from file %v", file))
	loadConfig(file)
	err := envconfig.Process("", config)
	if err != nil {
		return api_error.NewInternalServerError("Could not initalize configuration. Check your environment variables", err)
	}
	setDefaults(config)
	logger.Info("Done initalizing configuration")
	return nil
}

func setDefaults(config *AppConfig) {
}

func loadConfig(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		logger.Info("Could not open env file. Using Environment variable and defaults")
		return err
	}
	return nil
}
