"""Tact MCP Server - Main entry point."""

import asyncio
import json
from datetime import date
from typing import Any

from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import TextContent, Tool

from tact_mcp.client import TactClient

# Create the MCP server
server = Server("tact")

# Global client instance
_client: TactClient | None = None


def get_client() -> TactClient:
    """Get or create the Tact API client."""
    global _client
    if _client is None:
        _client = TactClient()
    return _client


def json_response(data: Any) -> list[TextContent]:
    """Format response as JSON text content."""
    return [TextContent(type="text", text=json.dumps(data, indent=2, default=str))]


def error_response(error: str) -> list[TextContent]:
    """Format error response."""
    return [TextContent(type="text", text=json.dumps({"error": error}, indent=2))]


@server.list_tools()
async def list_tools() -> list[Tool]:
    """List all available tools."""
    return [
        # Entry tools
        Tool(
            name="create_entry",
            description="Create a new time entry using natural language",
            inputSchema={
                "type": "object",
                "properties": {
                    "raw_text": {
                        "type": "string",
                        "description": "Natural language description (e.g., '2 hours development on Project Alpha')",
                    },
                    "entry_date": {
                        "type": "string",
                        "description": "Optional date in YYYY-MM-DD format",
                    },
                },
                "required": ["raw_text"],
            },
        ),
        Tool(
            name="list_entries",
            description="List time entries with optional filters",
            inputSchema={
                "type": "object",
                "properties": {
                    "status": {
                        "type": "string",
                        "enum": ["pending", "parsed", "needs_review", "failed"],
                        "description": "Filter by status",
                    },
                    "time_code_id": {
                        "type": "string",
                        "description": "Filter by time code ID",
                    },
                    "work_type_id": {
                        "type": "string",
                        "description": "Filter by work type ID",
                    },
                    "from_date": {
                        "type": "string",
                        "description": "Start date (YYYY-MM-DD)",
                    },
                    "to_date": {
                        "type": "string",
                        "description": "End date (YYYY-MM-DD)",
                    },
                    "limit": {
                        "type": "integer",
                        "description": "Maximum number of entries to return",
                        "default": 100,
                    },
                    "offset": {
                        "type": "integer",
                        "description": "Number of entries to skip",
                        "default": 0,
                    },
                },
            },
        ),
        Tool(
            name="get_entry",
            description="Get a single time entry by ID",
            inputSchema={
                "type": "object",
                "properties": {
                    "entry_id": {"type": "string", "description": "Entry ID (UUID)"},
                },
                "required": ["entry_id"],
            },
        ),
        Tool(
            name="update_entry",
            description="Update a time entry",
            inputSchema={
                "type": "object",
                "properties": {
                    "entry_id": {"type": "string", "description": "Entry ID (UUID)"},
                    "duration_minutes": {
                        "type": "integer",
                        "description": "Duration in minutes",
                    },
                    "time_code_id": {
                        "type": "string",
                        "description": "Time code ID",
                    },
                    "work_type_id": {
                        "type": "string",
                        "description": "Work type ID",
                    },
                    "description": {
                        "type": "string",
                        "description": "Entry description",
                    },
                    "entry_date": {
                        "type": "string",
                        "description": "Entry date (YYYY-MM-DD)",
                    },
                },
                "required": ["entry_id"],
            },
        ),
        Tool(
            name="delete_entry",
            description="Delete a time entry",
            inputSchema={
                "type": "object",
                "properties": {
                    "entry_id": {"type": "string", "description": "Entry ID (UUID)"},
                },
                "required": ["entry_id"],
            },
        ),
        Tool(
            name="reparse_entry",
            description="Re-trigger parsing for a time entry",
            inputSchema={
                "type": "object",
                "properties": {
                    "entry_id": {"type": "string", "description": "Entry ID (UUID)"},
                },
                "required": ["entry_id"],
            },
        ),
        # Time Code tools
        Tool(
            name="list_time_codes",
            description="List all time codes (billing categories)",
            inputSchema={
                "type": "object",
                "properties": {
                    "active_only": {
                        "type": "boolean",
                        "description": "Only return active time codes",
                    },
                    "project_id": {
                        "type": "string",
                        "description": "Filter by project ID",
                    },
                },
            },
        ),
        Tool(
            name="get_time_code",
            description="Get a single time code by ID",
            inputSchema={
                "type": "object",
                "properties": {
                    "time_code_id": {
                        "type": "string",
                        "description": "Time code ID",
                    },
                },
                "required": ["time_code_id"],
            },
        ),
        Tool(
            name="create_time_code",
            description="Create a new time code",
            inputSchema={
                "type": "object",
                "properties": {
                    "id": {
                        "type": "string",
                        "description": "Unique ID (e.g., 'PROJ-ALPHA')",
                    },
                    "name": {
                        "type": "string",
                        "description": "Display name",
                    },
                    "description": {
                        "type": "string",
                        "description": "Description for LLM context",
                    },
                    "keywords": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "Keywords for matching",
                    },
                    "examples": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "Example entries that map to this code",
                    },
                    "project_id": {
                        "type": "string",
                        "description": "Project ID (defaults to 'default')",
                        "default": "default",
                    },
                },
                "required": ["id", "name"],
            },
        ),
        Tool(
            name="update_time_code",
            description="Update a time code",
            inputSchema={
                "type": "object",
                "properties": {
                    "time_code_id": {
                        "type": "string",
                        "description": "Time code ID",
                    },
                    "name": {"type": "string", "description": "Display name"},
                    "description": {"type": "string", "description": "Description"},
                    "keywords": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "Keywords for matching",
                    },
                    "examples": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "Example entries",
                    },
                    "active": {"type": "boolean", "description": "Active status"},
                    "project_id": {
                        "type": "string",
                        "description": "Project ID to reassign time code",
                    },
                },
                "required": ["time_code_id"],
            },
        ),
        Tool(
            name="delete_time_code",
            description="Deactivate a time code (soft delete)",
            inputSchema={
                "type": "object",
                "properties": {
                    "time_code_id": {
                        "type": "string",
                        "description": "Time code ID",
                    },
                },
                "required": ["time_code_id"],
            },
        ),
        # Work Type tools
        Tool(
            name="list_work_types",
            description="List all work types",
            inputSchema={
                "type": "object",
                "properties": {
                    "active_only": {
                        "type": "boolean",
                        "description": "Only return active work types",
                    },
                },
            },
        ),
        Tool(
            name="get_work_type",
            description="Get a single work type by ID",
            inputSchema={
                "type": "object",
                "properties": {
                    "work_type_id": {
                        "type": "string",
                        "description": "Work type ID",
                    },
                },
                "required": ["work_type_id"],
            },
        ),
        Tool(
            name="create_work_type",
            description="Create a new work type",
            inputSchema={
                "type": "object",
                "properties": {
                    "name": {
                        "type": "string",
                        "description": "Work type name (e.g., 'Development')",
                    },
                },
                "required": ["name"],
            },
        ),
        Tool(
            name="update_work_type",
            description="Update a work type",
            inputSchema={
                "type": "object",
                "properties": {
                    "work_type_id": {
                        "type": "string",
                        "description": "Work type ID",
                    },
                    "name": {"type": "string", "description": "Display name"},
                    "active": {"type": "boolean", "description": "Active status"},
                },
                "required": ["work_type_id"],
            },
        ),
        Tool(
            name="delete_work_type",
            description="Deactivate a work type (soft delete)",
            inputSchema={
                "type": "object",
                "properties": {
                    "work_type_id": {
                        "type": "string",
                        "description": "Work type ID",
                    },
                },
                "required": ["work_type_id"],
            },
        ),
        # Project tools
        Tool(
            name="list_projects",
            description="List all projects",
            inputSchema={
                "type": "object",
                "properties": {
                    "active_only": {
                        "type": "boolean",
                        "description": "Only return active projects",
                    },
                },
            },
        ),
        Tool(
            name="get_project",
            description="Get a single project by ID",
            inputSchema={
                "type": "object",
                "properties": {
                    "project_id": {
                        "type": "string",
                        "description": "Project ID",
                    },
                },
                "required": ["project_id"],
            },
        ),
        Tool(
            name="create_project",
            description="Create a new project",
            inputSchema={
                "type": "object",
                "properties": {
                    "id": {
                        "type": "string",
                        "description": "Unique project ID (e.g., 'izg')",
                    },
                    "name": {
                        "type": "string",
                        "description": "Display name",
                    },
                },
                "required": ["id", "name"],
            },
        ),
        Tool(
            name="update_project",
            description="Update a project",
            inputSchema={
                "type": "object",
                "properties": {
                    "project_id": {
                        "type": "string",
                        "description": "Project ID",
                    },
                    "name": {"type": "string", "description": "Display name"},
                    "active": {"type": "boolean", "description": "Active status"},
                },
                "required": ["project_id"],
            },
        ),
        Tool(
            name="delete_project",
            description="Deactivate a project (soft delete)",
            inputSchema={
                "type": "object",
                "properties": {
                    "project_id": {
                        "type": "string",
                        "description": "Project ID",
                    },
                },
                "required": ["project_id"],
            },
        ),
        # Context tools
        Tool(
            name="list_project_context",
            description="List all context documents for a project",
            inputSchema={
                "type": "object",
                "properties": {
                    "project_id": {
                        "type": "string",
                        "description": "Project ID",
                    },
                },
                "required": ["project_id"],
            },
        ),
        Tool(
            name="list_time_code_context",
            description="List all context documents for a time code",
            inputSchema={
                "type": "object",
                "properties": {
                    "time_code_id": {
                        "type": "string",
                        "description": "Time code ID",
                    },
                },
                "required": ["time_code_id"],
            },
        ),
        Tool(
            name="create_project_context",
            description="Create a new context document for a project (used for RAG-enhanced parsing)",
            inputSchema={
                "type": "object",
                "properties": {
                    "project_id": {
                        "type": "string",
                        "description": "Project ID",
                    },
                    "content": {
                        "type": "string",
                        "description": "Context content text",
                    },
                },
                "required": ["project_id", "content"],
            },
        ),
        Tool(
            name="create_time_code_context",
            description="Create a new context document for a time code (used for RAG-enhanced parsing)",
            inputSchema={
                "type": "object",
                "properties": {
                    "time_code_id": {
                        "type": "string",
                        "description": "Time code ID",
                    },
                    "content": {
                        "type": "string",
                        "description": "Context content text",
                    },
                },
                "required": ["time_code_id", "content"],
            },
        ),
        Tool(
            name="get_context",
            description="Get a single context document by ID",
            inputSchema={
                "type": "object",
                "properties": {
                    "context_id": {
                        "type": "string",
                        "description": "Context document ID (UUID)",
                    },
                },
                "required": ["context_id"],
            },
        ),
        Tool(
            name="update_context",
            description="Update a context document (regenerates embedding)",
            inputSchema={
                "type": "object",
                "properties": {
                    "context_id": {
                        "type": "string",
                        "description": "Context document ID (UUID)",
                    },
                    "content": {
                        "type": "string",
                        "description": "New context content text",
                    },
                },
                "required": ["context_id", "content"],
            },
        ),
        Tool(
            name="delete_context",
            description="Delete a context document",
            inputSchema={
                "type": "object",
                "properties": {
                    "context_id": {
                        "type": "string",
                        "description": "Context document ID (UUID)",
                    },
                },
                "required": ["context_id"],
            },
        ),
    ]


@server.call_tool()
async def call_tool(name: str, arguments: dict[str, Any]) -> list[TextContent]:
    """Handle tool calls."""
    client = get_client()

    try:
        # Entry tools
        if name == "create_entry":
            entry_date = None
            if arguments.get("entry_date"):
                entry_date = date.fromisoformat(arguments["entry_date"])
            result = client.create_entry(arguments["raw_text"], entry_date)
            return json_response(result)

        elif name == "list_entries":
            from_date = None
            to_date = None
            if arguments.get("from_date"):
                from_date = date.fromisoformat(arguments["from_date"])
            if arguments.get("to_date"):
                to_date = date.fromisoformat(arguments["to_date"])
            result = client.list_entries(
                status=arguments.get("status"),
                time_code_id=arguments.get("time_code_id"),
                work_type_id=arguments.get("work_type_id"),
                from_date=from_date,
                to_date=to_date,
                limit=arguments.get("limit", 100),
                offset=arguments.get("offset", 0),
            )
            return json_response(result)

        elif name == "get_entry":
            result = client.get_entry(arguments["entry_id"])
            return json_response(result)

        elif name == "update_entry":
            entry_id = arguments["entry_id"]
            updates = {
                k: v for k, v in arguments.items() if k != "entry_id" and v is not None
            }
            result = client.update_entry(entry_id, **updates)
            return json_response(result)

        elif name == "delete_entry":
            result = client.delete_entry(arguments["entry_id"])
            return json_response(result)

        elif name == "reparse_entry":
            result = client.reparse_entry(arguments["entry_id"])
            return json_response(result)

        # Time Code tools
        elif name == "list_time_codes":
            result = client.list_time_codes(
                arguments.get("active_only"), arguments.get("project_id")
            )
            return json_response(result)

        elif name == "get_time_code":
            result = client.get_time_code(arguments["time_code_id"])
            return json_response(result)

        elif name == "create_time_code":
            result = client.create_time_code(
                id=arguments["id"],
                name=arguments["name"],
                description=arguments.get("description", ""),
                keywords=arguments.get("keywords"),
                examples=arguments.get("examples"),
                project_id=arguments.get("project_id", "default"),
            )
            return json_response(result)

        elif name == "update_time_code":
            time_code_id = arguments["time_code_id"]
            updates = {
                k: v
                for k, v in arguments.items()
                if k != "time_code_id" and v is not None
            }
            result = client.update_time_code(time_code_id, **updates)
            return json_response(result)

        elif name == "delete_time_code":
            result = client.delete_time_code(arguments["time_code_id"])
            return json_response(result)

        # Work Type tools
        elif name == "list_work_types":
            result = client.list_work_types(arguments.get("active_only"))
            return json_response(result)

        elif name == "get_work_type":
            result = client.get_work_type(arguments["work_type_id"])
            return json_response(result)

        elif name == "create_work_type":
            result = client.create_work_type(arguments["name"])
            return json_response(result)

        elif name == "update_work_type":
            work_type_id = arguments["work_type_id"]
            updates = {
                k: v
                for k, v in arguments.items()
                if k != "work_type_id" and v is not None
            }
            result = client.update_work_type(work_type_id, **updates)
            return json_response(result)

        elif name == "delete_work_type":
            result = client.delete_work_type(arguments["work_type_id"])
            return json_response(result)

        # Project tools
        elif name == "list_projects":
            result = client.list_projects(arguments.get("active_only"))
            return json_response(result)

        elif name == "get_project":
            result = client.get_project(arguments["project_id"])
            return json_response(result)

        elif name == "create_project":
            result = client.create_project(
                id=arguments["id"],
                name=arguments["name"],
            )
            return json_response(result)

        elif name == "update_project":
            project_id = arguments["project_id"]
            updates = {
                k: v
                for k, v in arguments.items()
                if k != "project_id" and v is not None
            }
            result = client.update_project(project_id, **updates)
            return json_response(result)

        elif name == "delete_project":
            result = client.delete_project(arguments["project_id"])
            return json_response(result)

        # Context tools
        elif name == "list_project_context":
            result = client.list_project_context(arguments["project_id"])
            return json_response(result)

        elif name == "list_time_code_context":
            result = client.list_time_code_context(arguments["time_code_id"])
            return json_response(result)

        elif name == "create_project_context":
            result = client.create_project_context(
                arguments["project_id"], arguments["content"]
            )
            return json_response(result)

        elif name == "create_time_code_context":
            result = client.create_time_code_context(
                arguments["time_code_id"], arguments["content"]
            )
            return json_response(result)

        elif name == "get_context":
            result = client.get_context(arguments["context_id"])
            return json_response(result)

        elif name == "update_context":
            result = client.update_context(
                arguments["context_id"], arguments["content"]
            )
            return json_response(result)

        elif name == "delete_context":
            result = client.delete_context(arguments["context_id"])
            return json_response(result)

        else:
            return error_response(f"Unknown tool: {name}")

    except Exception as e:
        return error_response(str(e))


async def run_server() -> None:
    """Run the MCP server with stdio transport."""
    async with stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            server.create_initialization_options(),
        )


def main() -> None:
    """Main entry point."""
    asyncio.run(run_server())


if __name__ == "__main__":
    main()
