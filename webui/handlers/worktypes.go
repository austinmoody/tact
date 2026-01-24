package handlers

import (
	"net/http"

	"tact-webui/api"
	"tact-webui/templates/pages"
)

type WorkTypeHandler struct {
	client *api.Client
}

func (h *WorkTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	workTypes, err := h.client.FetchWorkTypes()
	if err != nil {
		pages.Error("Error", "Failed to load work types: "+err.Error()).Render(r.Context(), w)
		return
	}

	pages.WorkTypes(workTypes).Render(r.Context(), w)
}

func (h *WorkTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	_, err := h.client.CreateWorkType(name)
	if err != nil {
		pages.ErrorPartial("Failed to create work type: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return updated list
	workTypes, _ := h.client.FetchWorkTypes()
	pages.WorkTypeList(workTypes).Render(r.Context(), w)
}

func (h *WorkTypeHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	workTypes, _ := h.client.FetchWorkTypes()

	for _, wt := range workTypes {
		if wt.ID == id {
			pages.WorkTypeEdit(wt).Render(r.Context(), w)
			return
		}
	}

	pages.ErrorPartial("Work type not found").Render(r.Context(), w)
}

func (h *WorkTypeHandler) Row(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	workTypes, _ := h.client.FetchWorkTypes()

	for _, wt := range workTypes {
		if wt.ID == id {
			pages.WorkTypeRow(wt).Render(r.Context(), w)
			return
		}
	}

	pages.ErrorPartial("Work type not found").Render(r.Context(), w)
}

func (h *WorkTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	wt, err := h.client.UpdateWorkType(id, &name)
	if err != nil {
		pages.ErrorPartial("Failed to update work type: " + err.Error()).Render(r.Context(), w)
		return
	}

	pages.WorkTypeRow(*wt).Render(r.Context(), w)
}

func (h *WorkTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.client.DeleteWorkType(id); err != nil {
		pages.ErrorPartial("Failed to delete work type: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return empty to remove the row
	w.WriteHeader(http.StatusOK)
}
