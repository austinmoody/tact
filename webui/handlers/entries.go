package handlers

import (
	"net/http"

	"tact-webui/api"
	"tact-webui/templates/components"
	"tact-webui/templates/pages"
)

type EntryHandler struct {
	client *api.Client
}

func (h *EntryHandler) Home(w http.ResponseWriter, r *http.Request) {
	entries, err := h.client.FetchEntries(100)
	if err != nil {
		pages.Error("Error", "Failed to load entries: "+err.Error()).Render(r.Context(), w)
		return
	}

	timeCodes, _ := h.client.FetchTimeCodes()
	workTypes, _ := h.client.FetchWorkTypes()

	pages.Home(entries, timeCodes, workTypes).Render(r.Context(), w)
}

func (h *EntryHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := api.EntryFilter{
		Limit:    100,
		Status:   r.URL.Query().Get("status"),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}

	entries, err := h.client.FetchEntriesFiltered(filter)
	if err != nil {
		pages.ErrorPartial("Failed to load entries: " + err.Error()).Render(r.Context(), w)
		return
	}

	timeCodes, _ := h.client.FetchTimeCodes()
	workTypes, _ := h.client.FetchWorkTypes()

	components.EntryList(entries, timeCodes, workTypes).Render(r.Context(), w)
}

func (h *EntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	userInput := r.FormValue("user_input")
	if userInput == "" {
		http.Error(w, "User input is required", http.StatusBadRequest)
		return
	}

	_, err := h.client.CreateEntry(userInput)
	if err != nil {
		pages.ErrorPartial("Failed to create entry: " + err.Error()).Render(r.Context(), w)
		return
	}

	// Return updated list
	h.List(w, r)
}

func (h *EntryHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entry, err := h.client.FetchEntry(id)
	if err != nil {
		pages.ErrorPartial("Failed to load entry: " + err.Error()).Render(r.Context(), w)
		return
	}

	timeCodes, _ := h.client.FetchTimeCodes()
	workTypes, _ := h.client.FetchWorkTypes()

	components.EntryDetail(*entry, timeCodes, workTypes).Render(r.Context(), w)
}

func (h *EntryHandler) Row(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entry, err := h.client.FetchEntry(id)
	if err != nil {
		pages.ErrorPartial("Failed to load entry: " + err.Error()).Render(r.Context(), w)
		return
	}

	timeCodes, _ := h.client.FetchTimeCodes()
	workTypes, _ := h.client.FetchWorkTypes()

	components.EntryRow(*entry, timeCodes, workTypes).Render(r.Context(), w)
}

func (h *EntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	update := api.EntryUpdate{}

	if v := r.FormValue("user_input"); v != "" {
		update.UserInput = &v
	}
	if v := r.FormValue("entry_date"); v != "" {
		update.EntryDate = &v
	}
	if v := r.FormValue("time_code_id"); v != "" {
		update.TimeCodeID = &v
	}
	if v := r.FormValue("work_type_id"); v != "" {
		update.WorkTypeID = &v
	}

	entry, err := h.client.UpdateEntry(id, update, true)
	if err != nil {
		pages.ErrorPartial("Failed to update entry: " + err.Error()).Render(r.Context(), w)
		return
	}

	timeCodes, _ := h.client.FetchTimeCodes()
	workTypes, _ := h.client.FetchWorkTypes()

	components.EntryRow(*entry, timeCodes, workTypes).Render(r.Context(), w)
}

func (h *EntryHandler) Reparse(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	entry, err := h.client.ReparseEntry(id)
	if err != nil {
		pages.ErrorPartial("Failed to reparse entry: " + err.Error()).Render(r.Context(), w)
		return
	}

	timeCodes, _ := h.client.FetchTimeCodes()
	workTypes, _ := h.client.FetchWorkTypes()

	components.EntryRow(*entry, timeCodes, workTypes).Render(r.Context(), w)
}
