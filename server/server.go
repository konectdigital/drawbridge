package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/konectdigital/drawbridge/config"
	"github.com/konectdigital/drawbridge/log"
	"github.com/konectdigital/drawbridge/plugin"
	"github.com/konectdigital/drawbridge/proxy"
	"github.com/konectdigital/drawbridge/utils"
	"github.com/konectdigital/muxinator"

	// Need to import all plugins to call their init functions
	_ "github.com/konectdigital/drawbridge/plugin/log"
	_ "github.com/konectdigital/drawbridge/plugin/retry"
)

func ListenAndServe(configuration *config.Configuration) error {
	router := muxinator.NewRouter()

	for _, api := range configuration.APIs {
		upstreamURL, err := url.Parse(api.UpstreamURL)
		if err != nil {
			return err
		}

		// Create a new proxy
		p := proxy.New(upstreamURL, api.AllowCrossOrigin)

		// Strip the prefix before passing to the proxy. Without this, the proxy will make a
		// request to upstreamUrl/prefix/path instead of upstreamUrl/path.
		prefix := utils.AddSlashes(api.Prefix)
		h := http.StripPrefix(prefix, p)

		// Every API gets the log plugin
		l := config.Plugin{
			Name:    "log",
			Enabled: true,
		}

		plugins := append([]config.Plugin{l}, api.Plugins...)
		m, err := getPluginMiddleware(plugins)
		if err != nil {
			return err
		}

		router.Handle("", prefix+"*", h, m...)
	}

	port := configuration.Port
	if port == 0 {
		port = 80
	}

	log.Printf("Listening on port %d\n", port)
	return router.ListenAndServe(fmt.Sprintf(":%d", port))
}

func getPluginMiddleware(plugins []config.Plugin) ([]muxinator.Middleware, error) {
	var middlewares []muxinator.Middleware

	for _, opts := range plugins {
		// Skip if plugin is disabled
		if !opts.Enabled {
			continue
		}

		p, err := plugin.Find(opts.Name)
		if err != nil {
			return nil, err
		}

		m, err := p.Middleware(opts.Config)
		if err != nil {
			return nil, err
		}

		middlewares = append(middlewares, m)
	}

	return middlewares, nil
}
