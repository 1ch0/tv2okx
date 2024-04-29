package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/fatih/color"
	"github.com/go-openapi/spec"
	"github.com/spf13/cobra"

	"github.com/1ch0/tv2okx/cmd/server/app/options"
	"github.com/1ch0/tv2okx/pkg/server"
	"github.com/1ch0/tv2okx/pkg/server/utils/log"
	"github.com/1ch0/tv2okx/pkg/utils"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "run",
		Long: `The API Server services REST operations and provides the frontend to the
cluster's shared state through which all other components interact.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return Run(s)
		},
		SilenceUsage: true,
	}

	buildSwaggerCmd := &cobra.Command{
		Use:   "build-swagger",
		Short: "Build swagger documentation of apiserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "docs/apidoc/latest-swagger.json"
			if len(args) > 0 {
				name = args[0]
			}
			func() {
				swagger, err := buildSwagger(s)
				if err != nil {
					log.Logger.Fatal(err.Error())
				}
				outData, err := json.MarshalIndent(swagger, "", "\t")
				if err != nil {
					log.Logger.Fatal(err.Error())
				}
				swaggerFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
				if err != nil {
					log.Logger.Fatal(err.Error())
				}
				defer func() {
					if err := swaggerFile.Close(); err != nil {
						log.Logger.Errorf("close swagger file failure %s", err.Error())
					}
				}()
				_, err = swaggerFile.Write(outData)
				if err != nil {
					log.Logger.Fatal(err.Error())
				}
				fmt.Println("build swagger config file success")
			}()
			return nil
		},
	}

	cmd.AddCommand(buildSwaggerCmd)

	return cmd
}

// Run runs the specified APIServer. This should never exit.
func Run(s *options.ServerRunOptions) error {
	// The server is not terminal, there is no color default.
	// Force set to false, this is useful for the dry-run API.
	color.NoColor = false

	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if s.GenericServerRunOptions.Server.PprofAddr != "" {
		go utils.EnablePprof(s.GenericServerRunOptions.Server.PprofAddr, errChan)
	}

	go func() {
		if err := run(ctx, s, errChan); err != nil {
			errChan <- fmt.Errorf("failed to run apiserver: %w", err)
		}
	}()
	var term = make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		log.Logger.Infof("Received SIGTERM, exiting gracefully...")
	case err := <-errChan:
		log.Logger.Errorf("Received an error: %s, exiting gracefully...", err.Error())
		return err
	}
	log.Logger.Infof("See you next time!")
	return nil
}

func run(ctx context.Context, s *options.ServerRunOptions, errChan chan error) error {
	server := server.New(*s.GenericServerRunOptions)
	return server.Run(ctx, errChan)
}

func buildSwagger(s *options.ServerRunOptions) (*spec.Swagger, error) {
	server := server.New(*s.GenericServerRunOptions)
	config, err := server.BuildRestfulConfig()
	if err != nil {
		return nil, err
	}
	return restfulspec.BuildSwagger(*config), nil
}
