"""HTTP client for the Tact API."""

import os
from datetime import date
from typing import Any

import httpx


class TactClient:
    """Client for interacting with the Tact API."""

    def __init__(self, base_url: str | None = None):
        self.base_url = base_url or os.getenv("TACT_API_URL", "http://localhost:2100")
        self._client = httpx.Client(base_url=self.base_url, timeout=30.0)

    def close(self) -> None:
        """Close the HTTP client."""
        self._client.close()

    def _handle_response(self, response: httpx.Response) -> dict[str, Any]:
        """Handle HTTP response, raising on errors."""
        if response.status_code == 204:
            return {"success": True}
        response.raise_for_status()
        return response.json()

    # --- Entries ---

    def create_entry(
        self, raw_text: str, entry_date: date | None = None
    ) -> dict[str, Any]:
        """Create a new time entry."""
        data: dict[str, Any] = {"raw_text": raw_text}
        if entry_date:
            data["entry_date"] = entry_date.isoformat()
        response = self._client.post("/entries", json=data)
        return self._handle_response(response)

    def list_entries(
        self,
        status: str | None = None,
        time_code_id: str | None = None,
        work_type_id: str | None = None,
        from_date: date | None = None,
        to_date: date | None = None,
        limit: int = 100,
        offset: int = 0,
    ) -> list[dict[str, Any]]:
        """List entries with optional filters."""
        params: dict[str, Any] = {"limit": limit, "offset": offset}
        if status:
            params["status"] = status
        if time_code_id:
            params["time_code_id"] = time_code_id
        if work_type_id:
            params["work_type_id"] = work_type_id
        if from_date:
            params["from_date"] = from_date.isoformat()
        if to_date:
            params["to_date"] = to_date.isoformat()
        response = self._client.get("/entries", params=params)
        return self._handle_response(response)

    def get_entry(self, entry_id: str) -> dict[str, Any]:
        """Get a single entry by ID."""
        response = self._client.get(f"/entries/{entry_id}")
        return self._handle_response(response)

    def update_entry(self, entry_id: str, **updates: Any) -> dict[str, Any]:
        """Update an entry."""
        response = self._client.patch(f"/entries/{entry_id}", json=updates)
        return self._handle_response(response)

    def delete_entry(self, entry_id: str) -> dict[str, Any]:
        """Delete an entry."""
        response = self._client.delete(f"/entries/{entry_id}")
        return self._handle_response(response)

    def reparse_entry(self, entry_id: str) -> dict[str, Any]:
        """Trigger re-parsing of an entry."""
        response = self._client.post(f"/entries/{entry_id}/reparse")
        return self._handle_response(response)

    # --- Time Codes ---

    def list_time_codes(
        self, active_only: bool | None = None, project_id: str | None = None
    ) -> list[dict[str, Any]]:
        """List all time codes."""
        params = {}
        if active_only is not None:
            params["active"] = active_only
        if project_id is not None:
            params["project_id"] = project_id
        response = self._client.get("/time-codes", params=params)
        return self._handle_response(response)

    def get_time_code(self, time_code_id: str) -> dict[str, Any]:
        """Get a single time code by ID."""
        response = self._client.get(f"/time-codes/{time_code_id}")
        return self._handle_response(response)

    def create_time_code(
        self,
        id: str,
        name: str,
        description: str = "",
        keywords: list[str] | None = None,
        examples: list[str] | None = None,
        project_id: str = "default",
    ) -> dict[str, Any]:
        """Create a new time code."""
        data = {
            "id": id,
            "name": name,
            "description": description,
            "keywords": keywords or [],
            "examples": examples or [],
            "project_id": project_id,
        }
        response = self._client.post("/time-codes", json=data)
        return self._handle_response(response)

    def update_time_code(self, time_code_id: str, **updates: Any) -> dict[str, Any]:
        """Update a time code."""
        response = self._client.put(f"/time-codes/{time_code_id}", json=updates)
        return self._handle_response(response)

    def delete_time_code(self, time_code_id: str) -> dict[str, Any]:
        """Deactivate a time code (soft delete)."""
        response = self._client.delete(f"/time-codes/{time_code_id}")
        return self._handle_response(response)

    # --- Work Types ---

    def list_work_types(self, active_only: bool | None = None) -> list[dict[str, Any]]:
        """List all work types."""
        params = {}
        if active_only is not None:
            params["active"] = active_only
        response = self._client.get("/work-types", params=params)
        return self._handle_response(response)

    def get_work_type(self, work_type_id: str) -> dict[str, Any]:
        """Get a single work type by ID."""
        response = self._client.get(f"/work-types/{work_type_id}")
        return self._handle_response(response)

    def create_work_type(self, name: str) -> dict[str, Any]:
        """Create a new work type."""
        response = self._client.post("/work-types", json={"name": name})
        return self._handle_response(response)

    def update_work_type(self, work_type_id: str, **updates: Any) -> dict[str, Any]:
        """Update a work type."""
        response = self._client.put(f"/work-types/{work_type_id}", json=updates)
        return self._handle_response(response)

    def delete_work_type(self, work_type_id: str) -> dict[str, Any]:
        """Deactivate a work type (soft delete)."""
        response = self._client.delete(f"/work-types/{work_type_id}")
        return self._handle_response(response)

    # --- Projects ---

    def list_projects(self, active_only: bool | None = None) -> list[dict[str, Any]]:
        """List all projects."""
        params = {}
        if active_only is not None:
            params["active"] = active_only
        response = self._client.get("/projects", params=params)
        return self._handle_response(response)

    def get_project(self, project_id: str) -> dict[str, Any]:
        """Get a single project by ID."""
        response = self._client.get(f"/projects/{project_id}")
        return self._handle_response(response)

    def create_project(self, id: str, name: str) -> dict[str, Any]:
        """Create a new project."""
        data = {"id": id, "name": name}
        response = self._client.post("/projects", json=data)
        return self._handle_response(response)

    def update_project(self, project_id: str, **updates: Any) -> dict[str, Any]:
        """Update a project."""
        response = self._client.put(f"/projects/{project_id}", json=updates)
        return self._handle_response(response)

    def delete_project(self, project_id: str) -> dict[str, Any]:
        """Deactivate a project (soft delete)."""
        response = self._client.delete(f"/projects/{project_id}")
        return self._handle_response(response)

    # --- Context Documents ---

    def list_project_context(self, project_id: str) -> list[dict[str, Any]]:
        """List all context documents for a project."""
        response = self._client.get(f"/projects/{project_id}/context")
        return self._handle_response(response)

    def create_project_context(
        self, project_id: str, content: str
    ) -> dict[str, Any]:
        """Create a new context document for a project."""
        response = self._client.post(
            f"/projects/{project_id}/context", json={"content": content}
        )
        return self._handle_response(response)

    def list_time_code_context(self, time_code_id: str) -> list[dict[str, Any]]:
        """List all context documents for a time code."""
        response = self._client.get(f"/time-codes/{time_code_id}/context")
        return self._handle_response(response)

    def create_time_code_context(
        self, time_code_id: str, content: str
    ) -> dict[str, Any]:
        """Create a new context document for a time code."""
        response = self._client.post(
            f"/time-codes/{time_code_id}/context", json={"content": content}
        )
        return self._handle_response(response)

    def get_context(self, context_id: str) -> dict[str, Any]:
        """Get a single context document by ID."""
        response = self._client.get(f"/context/{context_id}")
        return self._handle_response(response)

    def update_context(self, context_id: str, content: str) -> dict[str, Any]:
        """Update a context document."""
        response = self._client.put(f"/context/{context_id}", json={"content": content})
        return self._handle_response(response)

    def delete_context(self, context_id: str) -> dict[str, Any]:
        """Delete a context document."""
        response = self._client.delete(f"/context/{context_id}")
        return self._handle_response(response)
