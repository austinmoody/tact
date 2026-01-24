package handlers

import (
	"net/http"

	"tact-webui/api"
	"tact-webui/templates/components"
	"tact-webui/templates/pages"
)

type ContextHandler struct {
	client *api.Client
}

func (h *ContextHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// We need to find the context document
	// Try project contexts first, then time code contexts
	projects, _ := h.client.FetchProjects()
	for _, p := range projects {
		docs, _ := h.client.FetchProjectContext(p.ID)
		for _, doc := range docs {
			if doc.ID == id {
				components.ContextEdit(doc).Render(r.Context(), w)
				return
			}
		}
	}

	timeCodes, _ := h.client.FetchTimeCodes()
	for _, tc := range timeCodes {
		docs, _ := h.client.FetchTimeCodeContext(tc.ID)
		for _, doc := range docs {
			if doc.ID == id {
				components.ContextEdit(doc).Render(r.Context(), w)
				return
			}
		}
	}

	pages.ErrorPartial("Context not found").Render(r.Context(), w)
}

func (h *ContextHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	doc, err := h.client.UpdateContext(id, content)
	if err != nil {
		pages.ErrorPartial("Failed to update context: " + err.Error()).Render(r.Context(), w)
		return
	}

	components.ContextItem(*doc).Render(r.Context(), w)
}

func (h *ContextHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.client.DeleteContext(id); err != nil {
		pages.ErrorPartial("Failed to delete context: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return empty to remove the item
	w.WriteHeader(http.StatusOK)
}
