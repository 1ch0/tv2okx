package filters

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"

	"github.com/1ch0/tv2okx/pkg/server/utils"
)

func TestJSCache(t *testing.T) {
	chain := utils.NewFilterChain(loadJS, JSCache)
	res1 := httptest.NewRecorder()
	u, err := url.Parse("/test.js?v=1")
	assert.Equal(t, err, nil)
	chain.ProcessFilter(&http.Request{Method: "GET", URL: u}, res1)
	assert.Equal(t, res1.Code, 200)
	assert.Equal(t, res1.Body.String(), jsContent)
	assert.Equal(t, res1.Header().Get(HeaderHitCache), "false")
	assert.Equal(t, jsFileCache.Len(), 1)
	data, ok := jsFileCache.Get(u.String())
	assert.Equal(t, ok, true)
	assert.Equal(t, data.(*cacheData).data.String(), jsContent)

	res2 := httptest.NewRecorder()
	chain = utils.NewFilterChain(loadJS, JSCache)
	chain.ProcessFilter(&http.Request{Method: "GET", URL: u}, res2)
	assert.Equal(t, res2.Code, 200)
	assert.Equal(t, res2.Body.String(), jsContent)

	assert.Equal(t, res2.Header().Get(HeaderHitCache), "true")

	u2, err := url.Parse("/test.js?v=2")
	assert.Equal(t, err, nil)
	res3 := httptest.NewRecorder()
	chain = utils.NewFilterChain(loadJS, JSCache)
	chain.ProcessFilter(&http.Request{Method: "GET", URL: u2}, res3)
	assert.Equal(t, res3.Code, 200)
	assert.Equal(t, res3.Body.String(), jsContent)
	assert.Equal(t, res3.Header().Get(HeaderHitCache), "false")
	assert.Equal(t, jsFileCache.Len(), 2)

	u3, err := url.Parse("/test.js?v=3")
	assert.NoError(t, err)
	res4 := httptest.NewRecorder()
	chain = utils.NewFilterChain(loadJSNotOK, JSCache)
	chain.ProcessFilter(&http.Request{Method: "GET", URL: u3}, res4)
	assert.Equal(t, res4.Code, 304)
	assert.Equal(t, res4.Body.String(), "")
	assert.Equal(t, res4.Header().Get(HeaderHitCache), "false")
	assert.Equal(t, jsFileCache.Len(), 2)
}

var jsContent = "console.log(\"hello\")"

func loadJS(req *http.Request, res http.ResponseWriter) {
	fmt.Printf("miss cache,path:%s \n", req.URL.String())
	res.WriteHeader(200)
	res.Write([]byte(jsContent))
	res.Header().Add(restful.HEADER_ContentType, "application/javascript")
}

func loadJSNotOK(req *http.Request, res http.ResponseWriter) {
	fmt.Printf("return not 200,path:%s \n", req.URL.String())
	res.WriteHeader(304)
	res.Header().Add(restful.HEADER_ContentType, "application/javascript")
}
