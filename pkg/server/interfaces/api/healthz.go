package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	apis "github.com/1ch0/tv2okx/pkg/server/interfaces/api/dto/v1"
	"github.com/1ch0/tv2okx/pkg/server/utils/bcode"
)

type healthz struct {
}

// NewHealthz return healthz
func NewHealthz() Interface {
	return &healthz{}
}

// GetWebServiceRoute Get return healthz
func (u healthz) GetWebServiceRoute() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/healthz").Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("api for healthz webhook management")

	tags := []string{"healthz"}

	ws.Route(ws.GET("").To(u.handleHealthz).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("handle healthz check").
		Returns(200, "OK", apis.EmptyResponse{}).
		Returns(400, "Bad Request", bcode.Bcode{}).
		Writes(apis.EmptyResponse{}))

	return ws
}

func (u healthz) handleHealthz(req *restful.Request, res *restful.Response) {
	if err := res.WriteEntity(apis.EmptyResponse{}); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
}
