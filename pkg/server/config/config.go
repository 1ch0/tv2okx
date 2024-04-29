package config

import (
	"fmt"
	"github.com/1ch0/tv2okx/pkg/server/utils/log"
	"time"

	"github.com/spf13/viper"

	"github.com/spf13/pflag"

	"github.com/google/uuid"

	"github.com/1ch0/tv2okx/pkg/server/infrastructure/datastore"
)

// Config config for server
type Config struct {
	Server Server `yaml:"server"`

	// Datastore config
	Datastore datastore.Config `yaml:"datastore"`

	OKX OKX `yaml:"okx"`
}

type Server struct {
	Model string
	// api server bind address
	BindAddr string
	// monitor metric path
	MetricPath string
	// PprofAddr the address for pprof to use while exporting profiling results.
	PprofAddr string
	// LeaderConfig for leader election
	LeaderConfig leaderConfig `yaml:"leader_config"`
}

type OKX struct {
	APIKey     string
	APISecret  string
	PassPhrase string
}

type leaderConfig struct {
	ID       string
	LockName string
	Duration time.Duration
}

// ReadConfig config for server
func ReadConfig(path, name, configType string) *Config {
	config := &Config{}
	vip := viper.New()
	vip.AddConfigPath(path)
	vip.SetConfigName(name)
	vip.SetConfigType(configType)

	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := vip.Unmarshal(&config); err != nil {
		panic(err)
	}

	SetDefault(config)
	log.Logger.Infof("config: %+v", config)
	return config
}

func SetDefault(config *Config) {
	if config.Server.BindAddr == "" {
		config.Server.BindAddr = "0.0.0.0:8000"
	}
	if config.Server.MetricPath == "" {
		config.Server.MetricPath = "/metrics"
	}
	if config.Server.PprofAddr == "" {
		config.Server.PprofAddr = ""
	}

	if config.Server.LeaderConfig.ID == "" {
		config.Server.LeaderConfig = leaderConfig{
			ID:       uuid.New().String(),
			LockName: "apiserver-lock",
			Duration: time.Second * 5,
		}
	}

	if config.Datastore.Type == "" {
		config.Datastore.Type = "mongodb"
	}
	if config.Datastore.Database == "" {
		config.Datastore.Database = "go-restful-template"
	}
	switch config.Server.Model {
	case "debug", "local":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// NewConfig  returns a Config struct with default values
func NewConfig() *Config {
	return &Config{
		Server: Server{
			BindAddr:   "0.0.0.0:8000",
			MetricPath: "/metrics",
			LeaderConfig: leaderConfig{
				ID:       uuid.New().String(),
				LockName: "apiserver-lock",
				Duration: time.Second * 5,
			},
			PprofAddr: "",
		},
		Datastore: datastore.Config{
			Type:     "mongodb",
			Database: "go-restful-template",
			URL:      "",
		},
	}
}

// Validate validate generic server run options
func (s *Config) Validate() []error {
	var errs []error

	if s.Datastore.Type != "mongodb" {
		errs = append(errs, fmt.Errorf("not support datastore type %s", s.Datastore.Type))
	}

	return errs
}

// AddFlags adds flags to the specified FlagSet
func (s *Config) AddFlags(fs *pflag.FlagSet, c *Config) {
	fs.StringVar(&s.Server.BindAddr, "bind-addr", c.Server.BindAddr, "The bind address used to serve the http APIs.")
	fs.StringVar(&s.Server.MetricPath, "metrics-path", c.Server.MetricPath, "The path to expose the metrics.")
	fs.StringVar(&s.Datastore.Type, "datastore-type", c.Datastore.Type, "Metadata storage driver type, support mongodb")
	fs.StringVar(&s.Datastore.Database, "datastore-database", c.Datastore.Database, "Metadata storage database name, takes effect when the storage driver is mongodb.")
	fs.StringVar(&s.Datastore.URL, "datastore-url", c.Datastore.URL, "Metadata storage database url,takes effect when the storage driver is mongodb.")
	fs.StringVar(&s.Server.LeaderConfig.ID, "id", c.Server.LeaderConfig.ID, "the holder identity name")
	fs.StringVar(&s.Server.LeaderConfig.LockName, "lock-name", c.Server.LeaderConfig.LockName, "the lease lock resource name")
	fs.DurationVar(&s.Server.LeaderConfig.Duration, "duration", c.Server.LeaderConfig.Duration, "the lease lock resource name")
	fs.StringVar(&s.Server.PprofAddr, "pprof-addr", c.Server.PprofAddr, "The address for pprof to use while exporting profiling results. The default value is empty which means do not expose it. Set it to address like :6666 to expose it.")
}
