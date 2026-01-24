package handlers

import (
	"net/http"

	"tact-webui/api"
	"tact-webui/templates/components"
	"tact-webui/templates/pages"
)

type TimeCodeHandler struct {
	client *api.Client
}

func (h *TimeCodeHandler) List(w http.ResponseWriter, r *http.Request) {
	timeCodes, err := h.client.FetchTimeCodes()
	if err != nil {
		pages.Error("Error", "Failed to load time codes: "+err.Error()).Render(r.Context(), w)
		return
	}

	projects, _ := h.client.FetchProjects()
	pages.TimeCodes(timeCodes, projects).Render(r.Context(), w)
}

func (h *TimeCodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	projectID := r.FormValue("project_id")
	name := r.FormValue("name")

	if id == "" || projectID == "" || name == "" {
		http.Error(w, "ID, project, and name are required", http.StatusBadRequest)
		return
	}

	_, err := h.client.CreateTimeCode(id, projectID, name)
	if err != nil {
		pages.ErrorPartial("Failed to create time code: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return updated list
	timeCodes, _ := h.client.FetchTimeCodes()
	projects, _ := h.client.FetchProjects()
	pages.TimeCodeList(timeCodes, projects).Render(r.Context(), w)
}

func (h *TimeCodeHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	timeCodes, _ := h.client.FetchTimeCodes()
	projects, _ := h.client.FetchProjects()

	for _, tc := range timeCodes {
		if tc.ID == id {
			pages.TimeCodeEdit(tc, projects).Render(r.Context(), w)
			return
		}
	}

	pages.ErrorPartial("Time code not found").Render(r.Context(), w)
}

func (h *TimeCodeHandler) Row(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	timeCodes, _ := h.client.FetchTimeCodes()
	projects, _ := h.client.FetchProjects()

	for _, tc := range timeCodes {
		if tc.ID == id {
			pages.TimeCodeRow(tc, projects).Render(r.Context(), w)
			return
		}
	}

	pages.ErrorPartial("Time code not found").Render(r.Context(), w)
}

func (h *TimeCodeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	tc, err := h.client.UpdateTimeCode(id, nil, &name)
	if err != nil {
		pages.ErrorPartial("Failed to update time code: " + err.Error()).Render(r.Context(), w)
		return
	}

	projects, _ := h.client.FetchProjects()
	pages.TimeCodeRow(*tc, projects).Render(r.Context(), w)
}

func (h *TimeCodeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.client.DeleteTimeCode(id); err != nil {
		pages.ErrorPartial("Failed to delete time code: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return empty to remove the row
	w.WriteHeader(http.StatusOK)
}

func (h *TimeCodeHandler) Context(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	docs, err := h.client.FetchTimeCodeContext(id)
	if err != nil {
		pages.ErrorPartial("Failed to load context: " + err.Error()).Render(r.Context(), w)
		return
	}

	components.ContextList(docs, "timecode", id).Render(r.Context(), w)
}

func (h *TimeCodeHandler) CreateContext(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	_, err := h.client.CreateTimeCodeContext(id, content)
	if err != nil {
		pages.ErrorPartial("Failed to create context: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return updated list
	docs, _ := h.client.FetchTimeCodeContext(id)
	components.ContextList(docs, "timecode", id).Render(r.Context(), w)
}
