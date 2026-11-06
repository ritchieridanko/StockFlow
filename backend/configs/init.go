package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Client   Client   `mapstructure:"client"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	Tracer   Tracer   `mapstructure:"tracer"`
}

func Init(path string) (*Config, error) {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env != "dev" && env != "prod" {
		env = "dev"
	}
	if path == "" {
		path = "./configs"
	}

	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	cfg.App.Env = env
	cfg.Client.Addr = cfg.Client.Host + ":" + strconv.Itoa(cfg.Client.Port)
	cfg.Server.Addr = cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	cfg.Tracer.Addr = cfg.Tracer.Host + ":" + strconv.Itoa(cfg.Tracer.Port)
	cfg.Database.DSN = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	if env == "prod" {
		cfg.Client.Addr = "https://" + cfg.Client.Addr
	} else {
		cfg.Client.Addr = "http://" + cfg.Client.Addr
	}

	return &cfg, nil
}
