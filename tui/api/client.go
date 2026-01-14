package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tact-tui/model"
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

func (c *Client) FetchTimeCodes() ([]model.TimeCode, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/time-codes")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch time codes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var timeCodes []model.TimeCode
	if err := json.NewDecoder(resp.Body).Decode(&timeCodes); err != nil {
		return nil, fmt.Errorf("failed to decode time codes: %w", err)
	}

	return timeCodes, nil
}

func (c *Client) FetchWorkTypes() ([]model.WorkType, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/work-types")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch work types: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var workTypes []model.WorkType
	if err := json.NewDecoder(resp.Body).Decode(&workTypes); err != nil {
		return nil, fmt.Errorf("failed to decode work types: %w", err)
	}

	return workTypes, nil
}

// Entry methods

func (c *Client) FetchEntries(limit int) ([]model.Entry, error) {
	url := fmt.Sprintf("%s/entries?limit=%d", c.baseURL, limit)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var entries []model.Entry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("failed to decode entries: %w", err)
	}

	return entries, nil
}

func (c *Client) CreateEntry(rawText string) (*model.Entry, error) {
	body := map[string]string{"raw_text": rawText}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/entries", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var entry model.Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode entry: %w", err)
	}

	return &entry, nil
}

func (c *Client) ReparseEntry(id string) (*model.Entry, error) {
	url := fmt.Sprintf("%s/entries/%s/reparse", c.baseURL, id)
	resp, err := c.httpClient.Post(url, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to reparse entry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var entry model.Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode entry: %w", err)
	}

	return &entry, nil
}

// Time Code mutation methods

type TimeCodeCreate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	Examples    []string `json:"examples,omitempty"`
}

type TimeCodeUpdate struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	Examples    []string `json:"examples,omitempty"`
}

func (c *Client) CreateTimeCode(id, name, description string, keywords, examples []string) (*model.TimeCode, error) {
	body := TimeCodeCreate{
		ID:          id,
		Name:        name,
		Description: description,
		Keywords:    keywords,
		Examples:    examples,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/time-codes", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create time code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var tc model.TimeCode
	if err := json.NewDecoder(resp.Body).Decode(&tc); err != nil {
		return nil, fmt.Errorf("failed to decode time code: %w", err)
	}

	return &tc, nil
}

func (c *Client) UpdateTimeCode(id string, updates TimeCodeUpdate) (*model.TimeCode, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/time-codes/"+id, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update time code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var tc model.TimeCode
	if err := json.NewDecoder(resp.Body).Decode(&tc); err != nil {
		return nil, fmt.Errorf("failed to decode time code: %w", err)
	}

	return &tc, nil
}

func (c *Client) DeleteTimeCode(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/time-codes/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete time code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// Work Type mutation methods

type WorkTypeCreate struct {
	Name string `json:"name"`
}

type WorkTypeUpdate struct {
	Name *string `json:"name,omitempty"`
}

func (c *Client) CreateWorkType(name string) (*model.WorkType, error) {
	body := WorkTypeCreate{Name: name}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/work-types", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create work type: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wt model.WorkType
	if err := json.NewDecoder(resp.Body).Decode(&wt); err != nil {
		return nil, fmt.Errorf("failed to decode work type: %w", err)
	}

	return &wt, nil
}

func (c *Client) UpdateWorkType(id string, updates WorkTypeUpdate) (*model.WorkType, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/work-types/"+id, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update work type: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wt model.WorkType
	if err := json.NewDecoder(resp.Body).Decode(&wt); err != nil {
		return nil, fmt.Errorf("failed to decode work type: %w", err)
	}

	return &wt, nil
}

func (c *Client) DeleteWorkType(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/work-types/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete work type: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// Project methods

type ProjectCreate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type ProjectUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (c *Client) FetchProjects() ([]model.Project, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/projects")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var projects []model.Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode projects: %w", err)
	}

	return projects, nil
}

func (c *Client) CreateProject(id, name, description string) (*model.Project, error) {
	body := ProjectCreate{
		ID:          id,
		Name:        name,
		Description: description,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/projects", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var project model.Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode project: %w", err)
	}

	return &project, nil
}

func (c *Client) UpdateProject(id string, updates ProjectUpdate) (*model.Project, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/projects/"+id, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var project model.Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode project: %w", err)
	}

	return &project, nil
}

func (c *Client) DeleteProject(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/projects/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// Context methods

type ContextCreate struct {
	Content string `json:"content"`
}

type ContextUpdate struct {
	Content string `json:"content"`
}

func (c *Client) FetchProjectContext(projectID string) ([]model.ContextDocument, error) {
	url := fmt.Sprintf("%s/projects/%s/context", c.baseURL, projectID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project context: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var docs []model.ContextDocument
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, fmt.Errorf("failed to decode context: %w", err)
	}

	return docs, nil
}

func (c *Client) CreateProjectContext(projectID, content string) (*model.ContextDocument, error) {
	body := ContextCreate{Content: content}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/projects/%s/context", c.baseURL, projectID)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var doc model.ContextDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode context: %w", err)
	}

	return &doc, nil
}

func (c *Client) FetchTimeCodeContext(timeCodeID string) ([]model.ContextDocument, error) {
	url := fmt.Sprintf("%s/time-codes/%s/context", c.baseURL, timeCodeID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch time code context: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var docs []model.ContextDocument
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, fmt.Errorf("failed to decode context: %w", err)
	}

	return docs, nil
}

func (c *Client) CreateTimeCodeContext(timeCodeID, content string) (*model.ContextDocument, error) {
	body := ContextCreate{Content: content}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/time-codes/%s/context", c.baseURL, timeCodeID)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var doc model.ContextDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode context: %w", err)
	}

	return &doc, nil
}

func (c *Client) UpdateContext(contextID, content string) (*model.ContextDocument, error) {
	body := ContextUpdate{Content: content}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/context/"+contextID, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update context: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var doc model.ContextDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode context: %w", err)
	}

	return &doc, nil
}

func (c *Client) DeleteContext(contextID string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/context/"+contextID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete context: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
