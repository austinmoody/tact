# Simplify Time Code Model

## Why

The TimeCode model has fields that are redundant now that RAG context handles matching:

- **`description`** - Typically duplicates the `name` field
- **`keywords`** - RAG context does this better with learned examples
- **`examples`** - Never used in the LLM prompt (dead code)

With context documents learning from corrections, these static fields add complexity without value.

## What Changes

Remove three fields from the TimeCode model:
- `description` (Text, required) → removed
- `keywords` (JSON array) → removed
- `examples` (JSON array) → removed

The simplified model will have:
- `id` (primary key)
- `project_id` (foreign key)
- `name` (required)
- `active` (boolean)
- timestamps

The LLM prompt will show just: `- FEDS-163: Meeting Support`

RAG context continues to handle all "what goes where" matching logic.
