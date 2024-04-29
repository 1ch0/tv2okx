package container

import (
	"github.com/1ch0/tv2okx/pkg/server/utils/log"
	"time"

	"github.com/barnettZQG/inject"
)

// NewContainer new a IoC container
func NewContainer() *Container {
	return &Container{
		graph: inject.Graph{},
	}
}

// Container the IoC container
type Container struct {
	graph inject.Graph
}

// Provides provide some beans with default name
func (c *Container) Provides(beans ...interface{}) error {
	for _, bean := range beans {
		if err := c.graph.Provide(&inject.Object{Value: bean}); err != nil {
			return err
		}
	}
	return nil
}

// ProvideWithName provide the bean with name
func (c *Container) ProvideWithName(name string, bean interface{}) error {
	return c.graph.Provide(&inject.Object{Name: name, Value: bean})
}

// Populate populate dependency fields for all beans.
// this function must be called after providing all beans
func (c *Container) Populate() error {
	start := time.Now()
	defer func() {
		log.Logger.Infof("populate the bean container take time %s", time.Now().Sub(start))
	}()
	return c.graph.Populate()
}
