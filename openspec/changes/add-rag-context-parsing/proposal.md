# Proposal: Add RAG-Based Context for Time Code Parsing

## Why

The current time code matching system uses keywords and examples stored directly on each time code. This works for simple scenarios but fails for complex categorization rules like:

- **Override rules**: "ALL meetings with APHL go to Task 1, regardless of topic"
- **Technology-based routing**: "ALL UI work (node, react, nextjs) goes to Task 5"
- **Disambiguation**: "Security scans = Task 2, but security incidents = Task 1"
- **Acronym expansion**: "IZG = IZ Gateway, Xform = Transform"

Real-world time tracking often involves multiple projects with different complexity levels. The IZG project has 8 time codes with intricate categorization rules, while other projects might have just 1-2 simple codes. The system needs to handle both without requiring manual project selection at entry time.

## What Changes

### 1. Add Project Entity

Projects group related time codes and hold shared context (terminology, acronyms). Every time code belongs to a project.

```
Project:
  id: str (PK)
  name: str
  description: str | None
  active: bool
  created_at: datetime
  updated_at: datetime
```

### 2. Add Context Documents

Context documents store chunks of text that provide categorization guidance. They can be associated with either a project (shared context) or a specific time code (code-specific rules).

```
ContextDocument:
  id: uuid (PK)
  project_id: str | None (FK, for shared project context)
  time_code_id: str | None (FK, for code-specific context)
  content: str (the actual context text)
  embedding: blob | None (vector embedding for similarity search)
  created_at: datetime
  updated_at: datetime

  -- Constraint: exactly one of project_id or time_code_id must be set
```

### 3. Add RAG Infrastructure

- Vector embeddings for context documents (using sentence-transformers or similar)
- Similarity search to retrieve relevant context at parse time
- SQLite with sqlite-vec extension for vector storage (keeps it simple, no external DB)

### 4. Update Time Code Model

Add project association:

```
TimeCode:
  ... existing fields ...
  project_id: str (FK to Project, required)
```

### 5. Update LLM Parsing Flow

Current flow:
1. Build prompt with time code list (id, name, description, keywords)
2. LLM matches raw text to time code

New flow:
1. Embed the raw entry text
2. Retrieve relevant context chunks via vector similarity
3. Build prompt with:
   - Retrieved context chunks (with their project/time_code associations)
   - Time code list
4. LLM reasons over context to pick the right time code

### 6. Add API Endpoints

Projects:
- `POST /projects` - Create project
- `GET /projects` - List projects
- `GET /projects/{id}` - Get project
- `PUT /projects/{id}` - Update project
- `DELETE /projects/{id}` - Soft-delete project

Context Documents:
- `POST /projects/{id}/context` - Add context to project
- `POST /time-codes/{id}/context` - Add context to time code
- `GET /projects/{id}/context` - List project context docs
- `GET /time-codes/{id}/context` - List time code context docs
- `DELETE /context/{id}` - Delete context document

## Scope

**In Scope:**
- Project entity and API
- Context document entity and API
- Vector embeddings for context (local model, no external API)
- RAG retrieval during parsing
- Migration to associate existing time codes with projects

**Out of Scope:**
- Automatic project detection from entry text (projects are organizational, RAG handles matching)
- Context document versioning/history
- Bulk import of context documents
- TUI updates for projects and context management (separate change)

## Affected Components

- `backend/src/tact/db/models.py` - New Project and ContextDocument models
- `backend/src/tact/schemas/` - New schemas for projects and context
- `backend/src/tact/routes/` - New routes for projects and context
- `backend/src/tact/llm/prompts.py` - Updated to include RAG context
- `backend/src/tact/llm/provider.py` - Updated ParseContext
- New: `backend/src/tact/rag/` - RAG infrastructure (embeddings, retrieval)
- Database migration for new tables and time_code.project_id

## Dependencies

- sqlite-vec or similar for vector storage in SQLite
- sentence-transformers or similar for local embeddings (no external API dependency)
