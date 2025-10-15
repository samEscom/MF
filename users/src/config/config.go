package config

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server Server `mapstructure:"server"`
	App    App    `mapstructure:"app"`
	Db     DB     `mapstructure:"db"`
}

type BillsApi struct {
	URL string `mapstructure:"url" required:"true"`
}

type DB struct {
	Host     string `mapstructure:"host" required:"true"`
	Port     int    `mapstructure:"port" required:"true"`
	User     string `mapstructure:"user" required:"true"`
	Password string `mapstructure:"password" required:"true"`
	DB       string `mapstructure:"name" required:"true"`
	PoolSize int    `mapstructure:"pool_size" required:"false"`
}

type Server struct {
	Port     int    `required:"true" mapstructure:"port"`
	Host     string `required:"true" mapstructure:"host"`
	BasePath string `required:"true" mapstructure:"base_path"`
}

type App struct {
	ServiceName     string `mapstructure:"service_name" required:"true"`
	Postfix         string `mapstructure:"postfix" required:"true"`
	LoggerDebugMode bool   `mapstructure:"logger_debug_mode" required:"true"`
}

var configuration Configuration

func Config() Configuration {
	return configuration
}

func Environments() error {
	switch {
	case os.Getenv("ENVIRONMENT") == "local":
		if err := GenEnvsFromFile(".env.local"); err != nil {
			return fmt.Errorf("error reading config file %w", err)
		}
	case os.Getenv("ENVIRONMENT") == "dev":
		if err := GenEnvsFromFile(".env.dev"); err != nil {
			return fmt.Errorf("error reading config file %w", err)
		}
	case os.Getenv("ENVIRONMENT") == "prod":
		if err := GenEnvsFromFile(".env.prod"); err != nil {
			return fmt.Errorf("error reading config file %w", err)
		}
	}

	if errR := ReadConfigFromEnv(&configuration); errR != nil {
		return fmt.Errorf("error getting configurations from env %w", errR)
	}

	return nil
}

func GenEnvsFromFile(fileName string) error {
	file, err := os.Open(fmt.Sprintf("%s/%s", rootDir(), fileName))

	if err != nil {
		return err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		wrt := strings.Split(scanner.Text(), "=")

		//nolint:gomnd
		if len(wrt) < 2 {
			continue
		}

		if err := os.Setenv(wrt[0], wrt[1]); err != nil {
			return err
		}
	}

	return nil
}

func rootDir() string {
	dir, _ := os.Getwd()
	return dir
}

func ReadConfigFromEnv(config interface{}) error {
	bindEnvs(config)

	if errU := viper.Unmarshal(&config); errU != nil {
		return errU
	}

	return nil
}

func bindEnvs(config interface{}, parts ...string) {
	vc := reflect.ValueOf(config)

	if vc.Kind() == reflect.Ptr {
		vc = vc.Elem()
	}

	for i := 0; i < vc.NumField(); i++ {
		field := vc.Field(i)
		structField := vc.Type().Field(i)
		tv, ok := structField.Tag.Lookup("mapstructure")

		if !ok {
			continue
		}

		if field.Kind() == reflect.Struct {
			bindEnvs(field.Interface(), append(parts, tv)...)
		} else {
			_ = viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
