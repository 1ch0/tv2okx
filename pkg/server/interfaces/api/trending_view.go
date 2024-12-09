package api

import (
	"github.com/1ch0/tv2okx/pkg/server/utils"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/1ch0/tv2okx/pkg/server/domain/service"
	apis "github.com/1ch0/tv2okx/pkg/server/interfaces/api/dto/v1"
	"github.com/1ch0/tv2okx/pkg/server/utils/bcode"
)

type trendingView struct {
	TrendingView service.TrendingViewService `inject:""`
}

// NewTrendingView return trendingView
func NewTrendingView() Interface {
	return &trendingView{}
}

// GetWebServiceRoute Get return trendingView
func (u trendingView) GetWebServiceRoute() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/trendingView").Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("api for trendingView webhook management")

	tags := []string{"trendingView"}

	ws.Route(ws.POST("/webhook").To(u.handleWebhook).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc("handle trendingView webhook").
		Reads(apis.TrendingViewRequest{}).
		Returns(200, "OK", apis.EmptyResponse{}).
		Returns(400, "Bad Request", bcode.Bcode{}).
		Writes(apis.EmptyResponse{}))

	return ws
}

func (u trendingView) handleWebhook(req *restful.Request, res *restful.Response) {
	var tvReq *apis.TrendingViewRequest
	if err := req.ReadEntity(&tvReq); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
	if err := utils.Validate.Struct(tvReq); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}

	err := u.TrendingView.Webhook(req.Request.Context(), tvReq)
	if err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
	if err := res.WriteEntity(apis.EmptyResponse{}); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
}
