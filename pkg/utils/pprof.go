package utils

import (
	"encoding/json"
	"github.com/1ch0/tv2okx/pkg/server/utils/log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"
)

// EnablePprof listen to the pprofAddr and export the profiling results
// If the errChan is nil, this function will panic when the listening error occurred.
func EnablePprof(pprofAddr string, errChan chan error) {
	// Start pprof server if enabled
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/mem/stat", func(writer http.ResponseWriter, request *http.Request) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		bs, _ := json.Marshal(ms)
		_, _ = writer.Write(bs)
	})
	mux.HandleFunc("/gc", func(writer http.ResponseWriter, request *http.Request) {
		runtime.GC()
	})
	pprofServer := http.Server{
		Addr:              pprofAddr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Logger.Infof("Starting debug HTTP server", "addr", pprofServer.Addr)

	if err := pprofServer.ListenAndServe(); err != nil {
		log.Logger.Error(err, "Failed to start debug HTTP server")
		if errChan != nil {
			errChan <- err
		} else {
			panic(err)
		}
	}
}
