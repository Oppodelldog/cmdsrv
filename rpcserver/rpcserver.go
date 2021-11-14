package rpcserver

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/jsonrpc"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v3cdn"
)

func Run(opts ...Option) error {
	var (
		apiPath        = "/rpc"
		docsPathPrefix = "/docs"
		docsPath       = path.Join(docsPathPrefix, "openapi.json")

		c = applyOptions(defaultConfig(), opts...)
		h = Handler(c.Info)
		r = chi.NewRouter()
	)

	r.Mount(apiPath, h)
	r.Method(http.MethodGet, docsPath, h.OpenAPI)
	r.Mount(docsPathPrefix,
		v3cdn.NewHandlerWithConfig(
			swgui.Config{
				Title:       h.OpenAPI.Reflector().Spec.Info.Title,
				SwaggerJSON: docsPath,
				BasePath:    docsPathPrefix,
				SettingsUI:  jsonrpc.SwguiSettings(nil, apiPath),
			},
		),
	)

	for i := range c.Interactors {
		h.Add(c.Interactors[i])
	}

	return http.ListenAndServe(c.ServerAddr, r)
}
