package filters

import (
	"bytes"
	"github.com/1ch0/tv2okx/pkg/server/utils/log"
	"net/http"
	"strings"

	"github.com/1ch0/tv2okx/pkg/server/utils"
	"github.com/golang/groupcache/lru"
)

var jsFileCache = lru.New(100)

// HeaderHitCache the header key
var HeaderHitCache = "Hit-Cache"

func matchCacheCondition(req *http.Request) bool {
	return strings.HasSuffix(req.URL.Path, ".js") && req.Method == "GET"
}

// JSCache cache the JS static file.
func JSCache(req *http.Request, res http.ResponseWriter, chain *utils.FilterChain) {
	if matchCacheCondition(req) {
		if value, ok := jsFileCache.Get(req.URL.String()); ok {
			if cacheData, ok := value.(*cacheData); ok {
				if cacheData.data.Len() == 0 {
					log.Logger.Warnf("Cache data is empty")
					jsFileCache.Remove(req.URL.String())
				} else {
					cacheData.Write(res)
					return
				}
			}
		}
	}

	if matchCacheCondition(req) {
		res.Header().Set(HeaderHitCache, "false")
		cacheWriter := &CacheWriter{writer: res, cacheData: &cacheData{}}
		chain.ProcessFilter(req, cacheWriter)
		if cacheWriter.cacheData.code == http.StatusOK {
			jsFileCache.Add(req.URL.String(), cacheWriter.cacheData)
		} else {
			log.Logger.Warnf("Skip cache the js file, code: %d", cacheWriter.cacheData.code)
		}
		return
	}
	chain.ProcessFilter(req, res)
}

type cacheData struct {
	code   int
	data   bytes.Buffer
	header http.Header
}

func (c *cacheData) Write(w http.ResponseWriter) {
	for k, values := range c.header {
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}
	w.Header().Set(HeaderHitCache, "true")
	w.WriteHeader(c.code)
	if _, err := w.Write(c.data.Bytes()); err != nil {
		log.Logger.Errorf("failed to write the cache content, err: %s", err.Error())
	}
}

// CacheWriter generate the cache item the response body and status
type CacheWriter struct {
	writer    http.ResponseWriter
	cacheData *cacheData
}

// Header cache the header
func (c *CacheWriter) Header() http.Header {
	header := c.writer.Header()
	c.cacheData.header = header
	return header
}

// Write cache the data
func (c *CacheWriter) Write(b []byte) (int, error) {
	if _, err := c.cacheData.data.Write(b); err != nil {
		return -1, err
	}
	return c.writer.Write(b)
}

// WriteHeader cache the status code
func (c *CacheWriter) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
	c.cacheData.code = statusCode
}
