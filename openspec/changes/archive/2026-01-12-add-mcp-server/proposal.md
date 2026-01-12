# Change: Add MCP Server for Tact API

## Why

The Tact API is designed API-first to support multiple clients. An MCP (Model Context Protocol) server enables AI assistants like Claude Desktop, Goose, and GitHub Copilot to interact with Tact directly, turning any MCP-compatible client into a time tracking interface.

## What Changes

- Add new `mcp/` directory with Python MCP server
- Expose full Tact API as MCP tools (entries, time codes, work types, reports)
- Add Docker configuration for MCP server deployment
- Document client configuration for Claude Desktop, Goose, and other MCP clients

## Impact

- Affected specs: None (new capability)
- New spec: `mcp-server`
- Affected code:
  - New `mcp/` directory
  - `docker-compose.yml` (optional addition)
