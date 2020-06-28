package config

import (
	fileUtils "fantasymarket/utils/file"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is our global configuration file
type Config struct {
	Game        GameConfig     `mapstructure:"game"`
	TokenSecret string         `mapstructure:"tokenSecret"`
	LogLevel    string         `mapstructure:"logLevel"`
	Development bool           `mapstructure:"development"`
	Database    DatabaseConfig `mapstructure:"database"`
	Port        string         `mapstructure:"listenOn"`
}

// DatabaseConfig are specific settings for the database connection
type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // sqlite or postgres
	URL      string `mapstructure:"url"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSL      string `mapstructure:"ssl"`
}

// GameConfig are specific settings for the game mechanics
type GameConfig struct {
	TicksPerSecond  float64       `mapstructure:"ticksPerSecond"`  // How many times the game updates per second
	GameTimePerTick time.Duration `mapstructure:"gameTimePerTick"` // How much ingame time passes between updates
	StartDate       time.Time     `mapstructure:"startDate"`       // The initial ingame time
}

// Load loads the global configuration
func Load() (*Config, error) {

	// we populate a config.yaml in development environments
	_, configErr := os.Stat("config.yaml")
	_, exampleErr := os.Stat("config.example.yaml")

	// config.yaml is only created if the sample config exists in the same dir
	shouldCreateConfig := exampleErr == nil && os.IsNotExist(configErr)
	if shouldCreateConfig {
		if configErr = fileUtils.Copy("config.example.yaml", "config.yaml"); configErr != nil {
			return nil, configErr
		}
	}

	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.url", "")
	viper.SetDefault("database.host", "")
	viper.SetDefault("database.port", "")
	viper.SetDefault("database.username", "")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "")
	viper.SetDefault("database.ssl", "")

	viper.SetDefault("game.ticksPerSecond", 0.1)
	viper.SetDefault("game.timePerTick", time.Hour)

	viper.SetDefault("tokenSecret", "secret")
	viper.SetDefault("logLevel", "info")
	viper.SetDefault("development", "false")
	viper.SetDefault("port", "5000")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Find and read the config file
	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	conf := Config{
		Game: GameConfig{
			// viper somehow strugles with parsing dates so this is hardcoded
			StartDate: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	zerolog.SetGlobalLevel(getLogLevel(conf.LogLevel))

	// we have a nicer looking logger for development environments
	if conf.Development {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if !conf.Development && conf.TokenSecret == "secret" {
		log.Warn().Msg("tokenSecret should not be kept at the default value in production")
	}

	fmt.Println(conf)

	return &conf, err
}

func getLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.NoLevel
	}
}
