package filters

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"

	"github.com/1ch0/tv2okx/pkg/server/utils"
)

func TestGZip(t *testing.T) {
	chain := utils.NewFilterChain(loadJS, Gzip)
	res1 := httptest.NewRecorder()
	u, err := url.Parse("/test.js?v=1")
	assert.Equal(t, err, nil)
	reqHeader := http.Header{}
	reqHeader.Set(restful.HEADER_AcceptEncoding, restful.ENCODING_GZIP)
	chain.ProcessFilter(&http.Request{Method: "GET", URL: u, Header: reqHeader}, res1)
	assert.Equal(t, res1.Code, 200)
	assert.Equal(t, res1.HeaderMap.Get(restful.HEADER_ContentEncoding), restful.ENCODING_GZIP)

	// Gzip decode
	reader, err := gzip.NewReader(res1.Body)
	assert.Equal(t, err, nil)
	body, err := io.ReadAll(reader)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(body), jsContent)
}
