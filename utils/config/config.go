package config

import (
	"fantasymarket/utils/file"
	"fantasymarket/utils/timeutils"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uberConfig "go.uber.org/config"
)

// Config is our global configuration file
type Config struct {
	Game        GameConfig `yaml:"game"`
	TokenSecret string     `yaml:"tokenSecret"`
	LogLevel    string     `yaml:"logLevel"`
	Development bool       `yaml:"development"`
}

// GameConfig are specific settings for the game mechanics
type GameConfig struct {
	TicksPerSecond  float64        `yaml:"ticksPerSecond"`  // How many times the game updates per second
	GameTimePerTick time.Duration  `yaml:"gameTimePerTick"` // How much ingame time passes between updates
	StartDate       timeutils.Time `yaml:"startDate"`       // The initial ingame time}
}

var defaultConfig = Config{
	Game: GameConfig{
		TicksPerSecond:  0.1,
		StartDate:       timeutils.Time{Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
		GameTimePerTick: time.Hour,
	},
	TokenSecret: "secret",
	LogLevel:    "info",
}

// Load loads the global configuration
func Load() (*Config, error) {

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		err := file.Copy("config.example.yaml", "config.yaml")
		if err != nil {
			return nil, err
		}
	}

	yaml, err := uberConfig.NewYAML(
		uberConfig.Static(defaultConfig),
		uberConfig.File("config.yaml"),
		uberConfig.Expand(os.LookupEnv),
	)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := yaml.Get(uberConfig.Root).Populate(&conf); err != nil {
		return nil, err
	}

	zerolog.SetGlobalLevel(getLogLevel(conf.LogLevel))

	if conf.Development {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if !conf.Development && conf.TokenSecret == "secret" {
		log.Warn().Msg("tokenSecret should not be kept at the default value in production")
	}

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
