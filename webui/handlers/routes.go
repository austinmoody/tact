package handlers

import (
	"net/http"

	"tact-webui/api"
)

func RegisterRoutes(mux *http.ServeMux, client *api.Client) {
	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Entry handlers
	eh := &EntryHandler{client: client}
	mux.HandleFunc("GET /{$}", eh.Home)
	mux.HandleFunc("GET /entries/list", eh.List)
	mux.HandleFunc("POST /entries", eh.Create)
	mux.HandleFunc("GET /entries/{id}/detail", eh.Detail)
	mux.HandleFunc("GET /entries/{id}/row", eh.Row)
	mux.HandleFunc("PATCH /entries/{id}", eh.Update)
	mux.HandleFunc("POST /entries/{id}/reparse", eh.Reparse)

	// Timer handlers
	th := &TimerHandler{client: client}
	mux.HandleFunc("GET /timer", th.Page)
	mux.HandleFunc("GET /timer/stream", th.Stream)
	mux.HandleFunc("POST /timer/start", th.Start)
	mux.HandleFunc("POST /timer/pause", th.Pause)
	mux.HandleFunc("POST /timer/resume", th.Resume)
	mux.HandleFunc("POST /timer/stop", th.Stop)
	mux.HandleFunc("POST /timer/discard", th.Discard)

	// Project handlers
	ph := &ProjectHandler{client: client}
	mux.HandleFunc("GET /projects", ph.List)
	mux.HandleFunc("GET /projects/list", ph.ListPartial)
	mux.HandleFunc("POST /projects", ph.Create)
	mux.HandleFunc("GET /projects/{id}/edit", ph.Edit)
	mux.HandleFunc("GET /projects/{id}/row", ph.Row)
	mux.HandleFunc("PUT /projects/{id}", ph.Update)
	mux.HandleFunc("DELETE /projects/{id}", ph.Delete)
	mux.HandleFunc("GET /projects/{id}/context", ph.Context)
	mux.HandleFunc("POST /projects/{id}/context", ph.CreateContext)

	// Time Code handlers
	tch := &TimeCodeHandler{client: client}
	mux.HandleFunc("GET /time-codes", tch.List)
	mux.HandleFunc("POST /time-codes", tch.Create)
	mux.HandleFunc("GET /time-codes/{id}/edit", tch.Edit)
	mux.HandleFunc("GET /time-codes/{id}/row", tch.Row)
	mux.HandleFunc("PUT /time-codes/{id}", tch.Update)
	mux.HandleFunc("DELETE /time-codes/{id}", tch.Delete)
	mux.HandleFunc("GET /time-codes/{id}/context", tch.Context)
	mux.HandleFunc("POST /time-codes/{id}/context", tch.CreateContext)

	// Work Type handlers
	wth := &WorkTypeHandler{client: client}
	mux.HandleFunc("GET /work-types", wth.List)
	mux.HandleFunc("POST /work-types", wth.Create)
	mux.HandleFunc("GET /work-types/{id}/edit", wth.Edit)
	mux.HandleFunc("GET /work-types/{id}/row", wth.Row)
	mux.HandleFunc("PUT /work-types/{id}", wth.Update)
	mux.HandleFunc("DELETE /work-types/{id}", wth.Delete)

	// Context handlers (shared)
	ch := &ContextHandler{client: client}
	mux.HandleFunc("GET /context/{id}/edit", ch.Edit)
	mux.HandleFunc("PUT /context/{id}", ch.Update)
	mux.HandleFunc("DELETE /context/{id}", ch.Delete)
}
