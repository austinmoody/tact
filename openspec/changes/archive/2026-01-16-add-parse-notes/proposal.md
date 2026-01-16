# Proposal: Add Parse Notes Field

## Summary

Add a `parse_notes` field to time entries that captures the LLM's reasoning during parsing, including which context rules were considered and why a particular match was made (or not made).

## Motivation

Currently when an entry is marked `needs_review`, users don't know *why* the parser couldn't confidently match it. Was it:
- Ambiguous input matching multiple time codes?
- No relevant context documents?
- Conflicting rules?

Even for successfully parsed entries, understanding the reasoning helps users:
- Verify the parser made the right choice
- Learn how their context documents affect parsing
- Write clearer entries in the future

## Proposed Changes

### New Field: `parse_notes`

A text field on TimeEntry that contains:
1. **LLM reasoning** - A brief explanation of how the match was determined or why it couldn't be determined
2. **Context used** - Which RAG context document(s) were most relevant to the decision

Example for successful parse:
> "Matched to FEDS-163 based on APHL meeting rule. Used context: 'ALL meetings with APHL go to FEDS-163' (similarity: 0.85)"

Example for needs_review:
> "Could not confidently determine time code. Entry mentions both 'APHL' and 'UI work' which could map to different codes. Closest context: 'ALL UI work goes to FEDS-167' (similarity: 0.72)"

### Scope

**Backend:**
- Add `parse_notes` column to `time_entries` table (Alembic migration)
- Update `TimeEntry` model and `EntryResponse` schema
- Modify LLM prompt to request reasoning
- Update `ParseResult` to include notes
- Store closest RAG context info in notes

**TUI:**
- Display `parse_notes` in entry detail view (especially useful for `needs_review` entries)

**MCP:**
- No changes needed - field automatically included in API responses

## Benefits

- Users understand *why* parsing succeeded or failed
- Enables learning and improvement of context documents
- Debugging aid when entries are mis-categorized
- No significant performance impact (small additional LLM output)
