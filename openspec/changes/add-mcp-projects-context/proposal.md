# Proposal: Add MCP Tools for Projects and Context

## Summary

Add MCP tools to manage Projects and Context Documents, enabling AI clients to configure RAG-based parsing context through natural language.

## Motivation

The backend now supports Projects (for grouping time codes) and Context Documents (for RAG-enhanced parsing). MCP clients like Claude Desktop should be able to manage these entities to improve parsing accuracy.

## Scope

### In Scope

- Project management tools (list, get, create, update, delete)
- Context document tools for projects and time codes
- Client methods for all new API endpoints

### Out of Scope

- Bulk context import
- Context search/filtering beyond parent entity

## Design Decisions

1. **Consistent naming** - Follow existing patterns (`list_projects`, `create_project`, etc.)
2. **Separate context tools** - Use `list_project_context` and `list_time_code_context` rather than a generic tool with parent type parameter
3. **Generic context mutation** - Use `update_context` and `delete_context` that work on context ID regardless of parent type

## Affected Specs

- `mcp-server` - Add Project Tools and Context Tools requirements
