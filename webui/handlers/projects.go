package handlers

import (
	"net/http"
	"strings"

	"tact-webui/api"
	"tact-webui/model"
	"tact-webui/templates/components"
	"tact-webui/templates/pages"
)

type ProjectHandler struct {
	client *api.Client
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.client.FetchProjects()
	if err != nil {
		pages.Error("Error", "Failed to load projects: "+err.Error()).Render(r.Context(), w)
		return
	}

	pages.Projects(projects).Render(r.Context(), w)
}

func (h *ProjectHandler) ListPartial(w http.ResponseWriter, r *http.Request) {
	projects, err := h.client.FetchProjects()
	if err != nil {
		pages.ErrorPartial("Failed to load projects: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Filter by search term
	search := strings.ToLower(r.URL.Query().Get("search"))
	if search != "" {
		var filtered []model.Project
		for _, p := range projects {
			if strings.Contains(strings.ToLower(p.ID), search) ||
				strings.Contains(strings.ToLower(p.Name), search) {
				filtered = append(filtered, p)
			}
		}
		projects = filtered
	}

	pages.ProjectList(projects).Render(r.Context(), w)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")

	if id == "" || name == "" {
		http.Error(w, "ID and name are required", http.StatusBadRequest)
		return
	}

	_, err := h.client.CreateProject(id, name)
	if err != nil {
		pages.ErrorPartial("Failed to create project: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return updated list
	projects, _ := h.client.FetchProjects()
	pages.ProjectList(projects).Render(r.Context(), w)
}

func (h *ProjectHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projects, _ := h.client.FetchProjects()

	for _, p := range projects {
		if p.ID == id {
			pages.ProjectEdit(p).Render(r.Context(), w)
			return
		}
	}

	pages.ErrorPartial("Project not found").Render(r.Context(), w)
}

func (h *ProjectHandler) Row(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projects, _ := h.client.FetchProjects()

	for _, p := range projects {
		if p.ID == id {
			pages.ProjectRow(p).Render(r.Context(), w)
			return
		}
	}

	pages.ErrorPartial("Project not found").Render(r.Context(), w)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	project, err := h.client.UpdateProject(id, &name)
	if err != nil {
		pages.ErrorPartial("Failed to update project: " + err.Error()).Render(r.Context(), w)
		return
	}

	pages.ProjectRow(*project).Render(r.Context(), w)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.client.DeleteProject(id); err != nil {
		pages.ErrorPartial("Failed to delete project: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return empty to remove the row
	w.WriteHeader(http.StatusOK)
}

func (h *ProjectHandler) Context(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	docs, err := h.client.FetchProjectContext(id)
	if err != nil {
		pages.ErrorPartial("Failed to load context: " + err.Error()).Render(r.Context(), w)
		return
	}

	components.ContextList(docs, "project", id).Render(r.Context(), w)
}

func (h *ProjectHandler) CreateContext(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.client.CreateProjectContext(id, content)
	if err != nil {
		pages.ErrorPartial("Failed to create context: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return updated list
	docs, _ := h.client.FetchProjectContext(id)
	components.ContextList(docs, "project", id).Render(r.Context(), w)
}
