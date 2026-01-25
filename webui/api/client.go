package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"tact-webui/model"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTP helper methods

func (c *Client) get(path string, result interface{}) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	resp, err := c.httpClient.Post(c.baseURL+path, "application/json", reqBody)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) put(path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) patch(path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, c.baseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) delete(path string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// Entry methods

type EntryFilter struct {
	Limit    int
	Status   string
	DateFrom string
	DateTo   string
}

func (c *Client) FetchEntries(limit int) ([]model.Entry, error) {
	return c.FetchEntriesFiltered(EntryFilter{Limit: limit})
}

func (c *Client) FetchEntriesFiltered(filter EntryFilter) ([]model.Entry, error) {
	var entries []model.Entry
	path := fmt.Sprintf("/entries?limit=%d", filter.Limit)
	if filter.Status != "" {
		path += "&status=" + filter.Status
	}
	if filter.DateFrom != "" {
		path += "&from_date=" + filter.DateFrom
	}
	if filter.DateTo != "" {
		path += "&to_date=" + filter.DateTo
	}
	err := c.get(path, &entries)
	return entries, err
}

func (c *Client) FetchEntry(id string) (*model.Entry, error) {
	var entry model.Entry
	err := c.get("/entries/"+id, &entry)
	return &entry, err
}

func (c *Client) CreateEntry(userInput string) (*model.Entry, error) {
	body := map[string]string{"user_input": userInput}
	var entry model.Entry
	err := c.post("/entries", body, &entry)
	return &entry, err
}

type EntryUpdate struct {
	UserInput  *string `json:"user_input,omitempty"`
	EntryDate  *string `json:"entry_date,omitempty"`
	TimeCodeID *string `json:"time_code_id,omitempty"`
	WorkTypeID *string `json:"work_type_id,omitempty"`
}

func (c *Client) UpdateEntry(id string, update EntryUpdate, learn bool) (*model.Entry, error) {
	path := fmt.Sprintf("/entries/%s?learn=%t", id, learn)
	var entry model.Entry
	err := c.patch(path, update, &entry)
	return &entry, err
}

func (c *Client) ReparseEntry(id string) (*model.Entry, error) {
	var entry model.Entry
	err := c.post("/entries/"+id+"/reparse", nil, &entry)
	return &entry, err
}

// Project methods

func (c *Client) FetchProjects() ([]model.Project, error) {
	var projects []model.Project
	err := c.get("/projects", &projects)
	return projects, err
}

func (c *Client) CreateProject(id, name string) (*model.Project, error) {
	body := map[string]string{"id": id, "name": name}
	var project model.Project
	err := c.post("/projects", body, &project)
	return &project, err
}

func (c *Client) UpdateProject(id string, name *string) (*model.Project, error) {
	body := map[string]*string{"name": name}
	var project model.Project
	err := c.put("/projects/"+id, body, &project)
	return &project, err
}

func (c *Client) DeleteProject(id string) error {
	return c.delete("/projects/" + id)
}

// Time Code methods

func (c *Client) FetchTimeCodes() ([]model.TimeCode, error) {
	var timeCodes []model.TimeCode
	err := c.get("/time-codes", &timeCodes)
	return timeCodes, err
}

func (c *Client) CreateTimeCode(id, projectID, name string) (*model.TimeCode, error) {
	body := map[string]string{"id": id, "project_id": projectID, "name": name}
	var tc model.TimeCode
	err := c.post("/time-codes", body, &tc)
	return &tc, err
}

func (c *Client) UpdateTimeCode(id string, projectID, name *string) (*model.TimeCode, error) {
	body := map[string]*string{"project_id": projectID, "name": name}
	var tc model.TimeCode
	err := c.put("/time-codes/"+id, body, &tc)
	return &tc, err
}

func (c *Client) DeleteTimeCode(id string) error {
	return c.delete("/time-codes/" + id)
}

// Work Type methods

func (c *Client) FetchWorkTypes() ([]model.WorkType, error) {
	var workTypes []model.WorkType
	err := c.get("/work-types", &workTypes)
	return workTypes, err
}

func (c *Client) CreateWorkType(name string) (*model.WorkType, error) {
	body := map[string]string{"name": name}
	var wt model.WorkType
	err := c.post("/work-types", body, &wt)
	return &wt, err
}

func (c *Client) UpdateWorkType(id string, name *string) (*model.WorkType, error) {
	body := map[string]*string{"name": name}
	var wt model.WorkType
	err := c.put("/work-types/"+id, body, &wt)
	return &wt, err
}

func (c *Client) DeleteWorkType(id string) error {
	return c.delete("/work-types/" + id)
}

// Context methods

func (c *Client) FetchProjectContext(projectID string) ([]model.ContextDocument, error) {
	var docs []model.ContextDocument
	err := c.get("/projects/"+projectID+"/context", &docs)
	return docs, err
}

func (c *Client) CreateProjectContext(projectID, content string) (*model.ContextDocument, error) {
	body := map[string]string{"content": content}
	var doc model.ContextDocument
	err := c.post("/projects/"+projectID+"/context", body, &doc)
	return &doc, err
}

func (c *Client) FetchTimeCodeContext(timeCodeID string) ([]model.ContextDocument, error) {
	var docs []model.ContextDocument
	err := c.get("/time-codes/"+timeCodeID+"/context", &docs)
	return docs, err
}

func (c *Client) CreateTimeCodeContext(timeCodeID, content string) (*model.ContextDocument, error) {
	body := map[string]string{"content": content}
	var doc model.ContextDocument
	err := c.post("/time-codes/"+timeCodeID+"/context", body, &doc)
	return &doc, err
}

func (c *Client) UpdateContext(contextID, content string) (*model.ContextDocument, error) {
	body := map[string]string{"content": content}
	var doc model.ContextDocument
	err := c.put("/context/"+contextID, body, &doc)
	return &doc, err
}

func (c *Client) DeleteContext(contextID string) error {
	return c.delete("/context/" + contextID)
}
