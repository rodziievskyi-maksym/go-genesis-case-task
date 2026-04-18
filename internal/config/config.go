package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
	errInit  error
)

type Config struct {
	Host            string        `env:"HOST" default:"0.0.0.0" validate:"required"`
	Port            string        `env:"PORT" default:"8080" validate:"required"`
	Env             string        `env:"ENV" default:"development" validate:"required"`
	ScannerInterval time.Duration `env:"SCANNER_INTERVAL" default:"5m" validate:"required"`
	APIKey          string        `env:"API_KEY" validate:"required"`

	RedisHost     string        `env:"REDIS_HOST" default:"localhost" validate:"required"`
	RedisPort     string        `env:"REDIS_PORT" default:"6379" validate:"required"`
	RedisPass     string        `env:"REDIS_PASS" default:""`
	RedisDB       int           `env:"REDIS_DB" default:"0"`
	RedisCacheTTL time.Duration `env:"REDIS_CACHE_TTL" default:"10m" validate:"required"`

	SMTPHost string `env:"SMTP_HOST" validate:"required"`
	SMTPPort string `env:"SMTP_PORT" validate:"required"`
	SMTPUser string `env:"SMTP_USER" validate:"required"`
	SMTPPass string `env:"SMTP_PASS" validate:"required"`
	SMTPFrom string `env:"SMTP_FROM"`

	PostgresDSN string `env:"POSTGRES_DSN" validate:"required"`

	GitHubToken string `env:"GITHUB_TOKEN" validate:"required"`
}

var durationType = reflect.TypeOf(time.Duration(0))

func loadFromEnv(cfg *Config) error {
	value := reflect.ValueOf(cfg).Elem()
	typ := value.Type()
	errs := make([]error, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		raw, exists := os.LookupEnv(envKey)
		// Keep defaults when env var is empty.
		if !exists || raw == "" {
			continue
		}

		if err := setFieldValue(value.Field(i), raw); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", envKey, err))
		}
	}

	return errors.Join(errs...)
}

func setFieldValue(field reflect.Value, raw string) error {
	fieldType := field.Type()

	if fieldType == durationType {
		duration, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}

		field.SetInt(int64(duration))

		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(raw, 10, fieldType.Bits())
		if err != nil {
			return fmt.Errorf("invalid integer: %w", err)
		}

		field.SetInt(parsed)
	default:
		return fmt.Errorf("unsupported field type %s", fieldType)
	}

	return nil
}

func NewConfig(validate *validator.Validate, envPath ...string) error {
	if validate == nil {
		return errors.New("validator is nil")
	}

	once.Do(func() {
		slog.Info("Initializing configuration...")

		if err := godotenv.Load(envPath...); err != nil {
			errInit = errors.Join(errors.New("failed to load environment variables from .env file"), err)
			return
		}

		cfg := &Config{}

		if err := defaults.Set(cfg); err != nil {
			errInit = errors.Join(errors.New("failed to apply default configuration values"), err)
			return
		}

		if err := loadFromEnv(cfg); err != nil {
			errInit = errors.Join(errors.New("failed to parse environment variables"), err)
			return
		}

		if err := validate.Struct(cfg); err != nil {
			errInit = errors.Join(errors.New("configuration validation failed"), err)
			return
		}

		instance = cfg
		errInit = nil

		slog.Info("Configuration initialized")
	})

	if errInit != nil {
		return errInit
	}

	return nil
}

func Cfg() *Config {
	return instance
}

func IsProduction() bool {
	return instance.Env == "production"
}

func GetServerAddress() string {
	return fmt.Sprintf("%s:%s", instance.Host, instance.Port)
}

func GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", instance.RedisHost, instance.RedisPort)
}
