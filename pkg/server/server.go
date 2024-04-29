package server

import (
	"context"
	"fmt"
	"github.com/1ch0/tv2okx/pkg/server/utils/log"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"strings"
	"time"

	"github.com/1ch0/tv2okx/pkg/server/config"
	"github.com/1ch0/tv2okx/pkg/server/domain/service"
	"github.com/1ch0/tv2okx/pkg/server/infrastructure/datastore"
	"github.com/1ch0/tv2okx/pkg/server/infrastructure/datastore/mongodb"
	"github.com/1ch0/tv2okx/pkg/server/interfaces/api"
	"github.com/1ch0/tv2okx/pkg/server/utils"
	"github.com/1ch0/tv2okx/pkg/server/utils/container"
	"github.com/1ch0/tv2okx/pkg/server/utils/filters"
	pkgUtils "github.com/1ch0/tv2okx/pkg/utils"
	restfulSpec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/go-openapi/spec"
)

const (
	// SwaggerConfigRoutePath the path to request the swagger config
	SwaggerConfigRoutePath = "/debug/apidocs.json"

	// BuildPublicRoutePath the route prefix to request the build static files.
	BuildPublicRoutePath = "/public/build"

	// PluginPublicRoutePath the route prefix to request the plugin static files.
	PluginPublicRoutePath = "/public/plugins/"

	// PluginProxyRoutePath the route prefix to request the plugin backend server.
	PluginProxyRoutePath = "/proxy/plugins/"

	// DexRoutePath the route prefix to request the dex service
	DexRoutePath = "/dex"

	// BuildPublicPath the route prefix to request the build static files.
	BuildPublicPath = "public/build"
)

// APIServer interface for call api server
type APIServer interface {
	Run(context.Context, chan error) error
	BuildRestfulConfig() (*restfulSpec.Config, error)
}

// restServer rest server
type restServer struct {
	webContainer  *restful.Container
	beanContainer *container.Container
	cfg           config.Config
	dataStore     datastore.DataStore
}

// New create api server with config data
func New(cfg config.Config) (a APIServer) {
	s := &restServer{
		webContainer:  restful.NewContainer(),
		beanContainer: container.NewContainer(),
		cfg:           cfg,
	}
	return s
}

func (s *restServer) buildIoCContainer() (err error) {
	if err := s.beanContainer.ProvideWithName("config", s.cfg); err != nil {
		return fmt.Errorf("fail to provides the Config bean to the container: %w", err)
	}

	var ds datastore.DataStore
	switch s.cfg.Datastore.Type {
	case datastore.TypeMongoDB:
		ds, err = mongodb.New(context.Background(), s.cfg.Datastore)
		if err != nil {
			return fmt.Errorf("create mongodb datastore instance failure %w", err)
		}

	default:
		return fmt.Errorf("not support datastore type %s", s.cfg.Datastore.Type)
	}
	log.Logger.Infof("connect to datastore %s success", s.cfg.Datastore.Type)
	s.dataStore = ds
	if err := s.beanContainer.ProvideWithName("datastore", s.dataStore); err != nil {
		return fmt.Errorf("fail to provides the datastore bean to the container: %w", err)
	}

	// domain
	if err := s.beanContainer.Provides(service.InitServiceBean()...); err != nil {
		return fmt.Errorf("fail to provides the service bean to the container: %w", err)
	}

	// interfaces
	if err := s.beanContainer.Provides(api.InitAPIBean()...); err != nil {
		return fmt.Errorf("fail to provides the api bean to the container: %w", err)
	}

	if err := s.beanContainer.Populate(); err != nil {
		return fmt.Errorf("fail to populate the bean container: %w", err)
	}
	return nil
}

func (s *restServer) Run(ctx context.Context, errChan chan error) error {

	// build the Ioc Container
	if err := s.buildIoCContainer(); err != nil {
		return err
	}

	// init database
	if err := service.InitData(ctx); err != nil {
		return fmt.Errorf("fail to init database %w", err)
	}

	s.RegisterAPIRoute()

	return s.startHTTP(ctx)
}

// BuildRestfulConfig build the restful config
// This function will build the smallest set of beans
func (s *restServer) BuildRestfulConfig() (*restfulSpec.Config, error) {
	if err := s.buildIoCContainer(); err != nil {
		return nil, err
	}
	config := s.RegisterAPIRoute()
	return &config, nil
}

// RegisterAPIRoute register the API route
func (s *restServer) RegisterAPIRoute() restfulSpec.Config {
	/* **************************************************************  */
	/* *************       Open API Route Group     *****************  */
	/* **************************************************************  */
	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{},
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization", "RefreshToken"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		CookiesAllowed: true,
		Container:      s.webContainer}
	s.webContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	s.webContainer.Filter(s.webContainer.OPTIONSFilter)
	s.webContainer.Filter(s.OPTIONSFilter)

	// Add request log
	s.webContainer.Filter(s.requestLog)

	// Register all custom api
	for _, handler := range api.GetRegisteredAPI() {
		s.webContainer.Add(handler.GetWebServiceRoute())
	}

	config := restfulSpec.Config{
		WebServices:                   s.webContainer.RegisteredWebServices(), // you control what services are visible
		APIPath:                       SwaggerConfigRoutePath,
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	s.webContainer.Add(restfulSpec.NewOpenAPIService(config))
	return config
}

func (s *restServer) requestLog(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if req.HeaderParameter("Upgrade") == "websocket" && req.HeaderParameter("Connection") == "Upgrade" {
		chain.ProcessFilter(req, resp)
		return
	}
	start := time.Now()
	c := utils.NewResponseCapture(resp.ResponseWriter)
	resp.ResponseWriter = c
	chain.ProcessFilter(req, resp)
	takeTime := time.Since(start)
	log.Logger.With(
		"clientIP", pkgUtils.Sanitize(utils.ClientIP(req.Request)),
		"path", pkgUtils.Sanitize(req.Request.URL.Path),
		"method", req.Request.Method,
		"status", c.StatusCode(),
		"time", takeTime.String(),
		"responseSize", len(c.Bytes()),
	)
}

func (s *restServer) OPTIONSFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if req.Request.Method != "OPTIONS" {
		chain.ProcessFilter(req, resp)
		return
	}
	resp.AddHeader(restful.HEADER_AccessControlAllowCredentials, "true")
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "tv2okx api doc",
			Description: fmt.Sprintf("go-restful-template api doc, build time: %s", time.Now()),
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "ed.",
					Email: "neoed174@gmail.com",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache License 2.0",
					URL:  "",
				},
			},
			Version: "v1beta1",
		},
	}
}

func (s *restServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var staticFilters []utils.FilterFunction

	staticFilters = append(staticFilters, filters.Gzip)
	switch {
	case strings.HasPrefix(req.URL.Path, SwaggerConfigRoutePath):
		s.webContainer.ServeHTTP(res, req)
		return
	case strings.HasPrefix(req.URL.Path, BuildPublicRoutePath):
		utils.NewFilterChain(func(req *http.Request, res http.ResponseWriter) {
			s.staticFiles(res, req, "./")
		}, staticFilters...).ProcessFilter(req, res)
		return
	default:
		for _, pre := range api.GetAPIPrefix() {
			if strings.HasPrefix(req.URL.Path, pre) {
				s.webContainer.ServeHTTP(res, req)
				return
			}
		}
		// Rewrite to index.html, which means this route is handled by frontend.
		req.URL.Path = "/"
		utils.NewFilterChain(func(req *http.Request, res http.ResponseWriter) {
			s.staticFiles(res, req, BuildPublicPath)
		}, staticFilters...).ProcessFilter(req, res)
	}
}

func (s *restServer) staticFiles(res http.ResponseWriter, req *http.Request, root string) {
	http.FileServer(http.Dir(root)).ServeHTTP(res, req)
}

func (s *restServer) startHTTP(ctx context.Context) error {
	// Start HTTP apiserver
	log.Logger.Infof("HTTP APIs are being served on: %s, ctx: %s", s.cfg.Server.BindAddr, ctx)
	server := &http.Server{Addr: s.cfg.Server.BindAddr, Handler: s.webContainer, ReadHeaderTimeout: 2 * time.Second}
	return server.ListenAndServe()
}
