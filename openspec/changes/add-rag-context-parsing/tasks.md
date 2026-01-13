# Tasks

## 1. Add Project Entity and API

- [ ] 1.1 Add Project model to `db/models.py`
- [ ] 1.2 Add Project schemas (`ProjectCreate`, `ProjectUpdate`, `ProjectResponse`)
- [ ] 1.3 Add Project routes (`POST`, `GET`, `GET/{id}`, `PUT/{id}`, `DELETE/{id}`)
- [ ] 1.4 Add tests for Project API endpoints

## 2. Update Time Code with Project Association

- [ ] 2.1 Add `project_id` foreign key to TimeCode model
- [ ] 2.2 Update TimeCode schemas to include `project_id`
- [ ] 2.3 Update TimeCode create/update to require `project_id`
- [ ] 2.4 Create database migration (add column, create default project, backfill)
- [ ] 2.5 Update TimeCode API tests

## 3. Add Context Document Entity and API

- [ ] 3.1 Add ContextDocument model to `db/models.py`
- [ ] 3.2 Add ContextDocument schemas
- [ ] 3.3 Add routes for project context (`POST /projects/{id}/context`, `GET /projects/{id}/context`)
- [ ] 3.4 Add routes for time code context (`POST /time-codes/{id}/context`, `GET /time-codes/{id}/context`)
- [ ] 3.5 Add route to delete context (`DELETE /context/{id}`)
- [ ] 3.6 Add tests for Context API endpoints

## 4. Add RAG Infrastructure

- [ ] 4.1 Add sqlite-vec dependency and setup
- [ ] 4.2 Add sentence-transformers dependency
- [ ] 4.3 Create `rag/` module with embedding functions
- [ ] 4.4 Create vector table/index for context embeddings
- [ ] 4.5 Implement `embed_and_store()` for context documents
- [ ] 4.6 Implement `retrieve_similar()` for query-time retrieval
- [ ] 4.7 Add tests for RAG retrieval

## 5. Integrate RAG with Context Document API

- [ ] 5.1 Auto-embed content when context document is created
- [ ] 5.2 Re-embed content when context document is updated
- [ ] 5.3 Remove embedding when context document is deleted

## 6. Update LLM Parsing to Use RAG

- [ ] 6.1 Update `ParseContext` to include retrieved context chunks
- [ ] 6.2 Update `prompts.py` to format retrieved context in prompt
- [ ] 6.3 Update parser to embed entry text and retrieve context before calling LLM
- [ ] 6.4 Add tests for RAG-enhanced parsing

## Dependencies

- Task 2 depends on Task 1 (need Project before adding FK)
- Task 3 depends on Task 1 and Task 2 (context references both)
- Task 5 depends on Task 3 and Task 4 (embedding requires both)
- Task 6 depends on Task 4 and Task 5 (parsing uses RAG)

## Notes

- Tasks 1-3 can be done incrementally and tested independently
- Task 4 (RAG infrastructure) can be developed in parallel with Tasks 1-3
- TUI updates (projects screen, context management) deferred to separate change
