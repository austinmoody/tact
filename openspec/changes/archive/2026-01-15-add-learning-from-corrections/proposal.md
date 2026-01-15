# Proposal: Add Learning from Manual Corrections

## Summary

When a user manually corrects a time entry, automatically create a context document for the associated time code so the system "learns" from the correction and improves future parsing.

## Motivation

Currently, when the LLM parses an entry with low confidence or incorrectly, users must manually correct it. However, these corrections don't feed back into the system - the same mistake may recur on similar entries. Users have to manually add context documents to improve parsing, which is tedious.

By automatically creating context documents from corrections, the system learns from user feedback without requiring explicit context management.

## Scope

### In Scope

- Add `learn` query parameter to PATCH `/entries/{id}` endpoint (default: `true`)
- When `learn=true` and the entry has a `time_code_id`, create a context document
- Format the context document content to be useful for future parsing
- Skip learning if no time_code_id is set on the entry

### Out of Scope

- Differentiating learned vs manually-created context documents
- Learning from work type corrections (work types don't have context yet)
- Deduplication of similar learned examples
- Automatic pruning/cleanup of learned context
- TUI changes (no UI needed - this is API-only)

## Design Decisions

1. **Query parameter over body field** - Using `?learn=false` keeps the request body clean and makes it easy to disable learning without changing the payload structure.

2. **Default to learning** - Most corrections are intentional improvements, so learning should be opt-out rather than opt-in.

3. **Time code context only** - Context documents attach to time codes (not work types) since that's where RAG context currently lives.

4. **Simple content format** - The learned context uses a human-readable format:
   ```
   Example: "2h standup and sprint planning"
   Parsed as: 120 minutes, work_type: meetings
   ```

5. **No differentiation** - Learned context documents are regular context documents. Users can view/delete them through existing context management.

## Affected Specs

- `backend-api` - Add learning behavior to entry update endpoint
