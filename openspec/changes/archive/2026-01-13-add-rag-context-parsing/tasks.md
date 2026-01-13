# Tasks

## 1. Add Project Entity and API

- [x] 1.1 Add Project model to `db/models.py`
- [x] 1.2 Add Project schemas (`ProjectCreate`, `ProjectUpdate`, `ProjectResponse`)
- [x] 1.3 Add Project routes (`POST`, `GET`, `GET/{id}`, `PUT/{id}`, `DELETE/{id}`)
- [x] 1.4 Add tests for Project API endpoints

## 2. Update Time Code with Project Association

- [x] 2.1 Add `project_id` foreign key to TimeCode model
- [x] 2.2 Update TimeCode schemas to include `project_id`
- [x] 2.3 Update TimeCode create/update to require `project_id`
- [x] 2.4 Create database migration (add column, create default project, backfill)
- [x] 2.5 Update TimeCode API tests

## 3. Add Context Document Entity and API

- [x] 3.1 Add ContextDocument model to `db/models.py`
- [x] 3.2 Add ContextDocument schemas
- [x] 3.3 Add routes for project context (`POST /projects/{id}/context`, `GET /projects/{id}/context`)
- [x] 3.4 Add routes for time code context (`POST /time-codes/{id}/context`, `GET /time-codes/{id}/context`)
- [x] 3.5 Add route to delete context (`DELETE /context/{id}`)
- [x] 3.6 Add tests for Context API endpoints

## 4. Add RAG Infrastructure

- [x] 4.1 Add sqlite-vec dependency and setup
- [x] 4.2 Add sentence-transformers dependency
- [x] 4.3 Create `rag/` module with embedding functions
- [x] 4.4 Create vector table/index for context embeddings
- [x] 4.5 Implement `embed_and_store()` for context documents
- [x] 4.6 Implement `retrieve_similar()` for query-time retrieval
- [x] 4.7 Add tests for RAG retrieval

## 5. Integrate RAG with Context Document API

- [x] 5.1 Auto-embed content when context document is created
- [x] 5.2 Re-embed content when context document is updated
- [x] 5.3 Remove embedding when context document is deleted

## 6. Update LLM Parsing to Use RAG

- [x] 6.1 Update `ParseContext` to include retrieved context chunks
- [x] 6.2 Update `prompts.py` to format retrieved context in prompt
- [x] 6.3 Update parser to embed entry text and retrieve context before calling LLM
- [x] 6.4 Add tests for RAG-enhanced parsing

## Dependencies

- Task 2 depends on Task 1 (need Project before adding FK)
- Task 3 depends on Task 1 and Task 2 (context references both)
- Task 5 depends on Task 3 and Task 4 (embedding requires both)
- Task 6 depends on Task 4 and Task 5 (parsing uses RAG)

## Notes

- Tasks 1-3 can be done incrementally and tested independently
- Task 4 (RAG infrastructure) can be developed in parallel with Tasks 1-3
- TUI updates (projects screen, context management) deferred to separate change
