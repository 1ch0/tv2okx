package options

import (
	"github.com/1ch0/tv2okx/pkg/server/config"
)

// ServerRunOptions contains everything necessary to create and run api server
type ServerRunOptions struct {
	GenericServerRunOptions *config.Config
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters
func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		GenericServerRunOptions: config.ReadConfig("./config", "config", "yaml"),
	}
	return s
}
