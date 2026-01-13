## Context

TACT currently uses a simple keyword/example-based system for time code matching. The LLM receives a list of time codes with their keywords and picks the best match. This fails for complex categorization scenarios that require:

1. **Override rules** - "APHL meetings always go to Task 1, even if they discuss UI"
2. **Contextual disambiguation** - "Security scans = Task 2, security incidents = Task 1"
3. **Technology routing** - "All React/Node work = Task 5"
4. **Shared terminology** - "IZG = IZ Gateway" applies to multiple codes

### Constraints

- Must work with SQLite (no external vector DB)
- Must use local embeddings (no external API dependency for embeddings)
- Must not require manual project selection when entering time
- Must support varying complexity (1 simple code vs 8 complex codes with rules)

## Goals / Non-Goals

**Goals:**
- Enable rich categorization rules via retrievable context
- Support project-level shared context (acronyms, terminology)
- Support time-code-specific context (rules for that code)
- Automatic retrieval of relevant context at parse time
- Simple API for managing context documents

**Non-Goals:**
- Automatic project detection as a separate step (RAG handles this implicitly)
- Context versioning or history tracking
- Real-time context updates during parsing (batch embedding is fine)
- Sophisticated chunking strategies (simple document-per-chunk for now)

## Decisions

### Project as Organizational Grouping

**Decision:** Every time code belongs to a project. Projects hold shared context.

**Rationale:**
- Keeps IZG acronyms in one place, not duplicated across 8 codes
- Adding a new IZG code automatically benefits from shared context
- Simple projects (1 code) still follow the same pattern for consistency
- Project is organizational only - not selected at entry time

**Data Model:**
```python
class Project(Base):
    __tablename__ = "projects"

    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(nullable=False)
    description: Mapped[str | None] = mapped_column(Text, default=None)
    active: Mapped[bool] = mapped_column(default=True)
    created_at: Mapped[datetime] = mapped_column(default=utc_now)
    updated_at: Mapped[datetime] = mapped_column(default=utc_now, onupdate=utc_now)
```

### Context Documents with Embeddings

**Decision:** Store context as documents with vector embeddings. Each document is associated with either a project or a time code.

**Rationale:**
- Documents are the natural unit for RAG
- Project-level docs hold shared context (acronyms, terminology)
- Time-code-level docs hold specific rules ("APHL meetings = this code")
- Vector similarity finds relevant context regardless of project/code

**Data Model:**
```python
class ContextDocument(Base):
    __tablename__ = "context_documents"

    id: Mapped[str] = mapped_column(primary_key=True, default=lambda: str(uuid.uuid4()))
    project_id: Mapped[str | None] = mapped_column(ForeignKey("projects.id"), default=None)
    time_code_id: Mapped[str | None] = mapped_column(ForeignKey("time_codes.id"), default=None)
    content: Mapped[str] = mapped_column(Text, nullable=False)
    embedding: Mapped[bytes | None] = mapped_column(LargeBinary, default=None)
    created_at: Mapped[datetime] = mapped_column(default=utc_now)
    updated_at: Mapped[datetime] = mapped_column(default=utc_now, onupdate=utc_now)

    # Constraint: exactly one of project_id or time_code_id must be set
```

### Vector Storage in SQLite

**Decision:** Use sqlite-vec extension for vector similarity search.

**Rationale:**
- Keeps everything in one SQLite file (portable, simple)
- No external vector DB dependency
- Good enough performance for personal use (hundreds/thousands of docs)
- Well-maintained extension with Python bindings

**Alternative considered:** Chroma or FAISS
- More features but adds complexity
- External files/processes to manage
- Overkill for this use case

### Local Embeddings

**Decision:** Use sentence-transformers with a small model (all-MiniLM-L6-v2 or similar).

**Rationale:**
- No external API calls for embeddings
- Fast enough for real-time embedding of entry text
- Good quality for semantic similarity
- Model downloads once, runs locally

**Implementation:**
```python
from sentence_transformers import SentenceTransformer

model = SentenceTransformer('all-MiniLM-L6-v2')

def embed_text(text: str) -> list[float]:
    return model.encode(text).tolist()
```

### RAG Retrieval Flow

**Decision:** Embed entry text, retrieve top-k similar context docs, include in LLM prompt.

**Flow:**
1. Entry submitted: "2h APHL meeting about UI deployment"
2. Embed entry text → vector
3. Query sqlite-vec for top-k similar context docs (k=5-10)
4. Build prompt with:
   - Retrieved context (each chunk labeled with its project/time_code source)
   - Full list of active time codes
5. LLM reasons: "Context says APHL meetings always go to FEDS-163"
6. Return: time_code_id=FEDS-163, duration=120, etc.

**Prompt Structure:**
```
Retrieved Context:
[Project: IZG] IZG = IZ Gateway. APHL = Association of Public Health Laboratories.
[FEDS-163] ALL meetings with APHL go to this code, regardless of meeting topic.
[FEDS-167] ALL UI-related work (node, react, nextjs) goes to this code.

Available Time Codes:
- FEDS-163: IZG (Task 1) System Ops
- FEDS-167: IZG (Task 5) Provider Onboarding
- ...

Parse this entry: "2h APHL meeting about UI deployment"
```

### Embedding on Create/Update

**Decision:** Embed context documents when they are created or updated, not at query time.

**Rationale:**
- Embedding is relatively slow (~50-100ms per doc)
- Context docs change rarely, entries are created frequently
- Pre-computed embeddings make query-time fast

**Implementation:**
- POST /context → embed content → store doc with embedding
- PUT /context → re-embed content → update doc

### Migration Strategy

**Decision:** Create a "default" project for existing time codes.

**Rationale:**
- Existing time codes must belong to a project (FK constraint)
- User can reorganize later (move codes to proper projects)
- Doesn't break existing functionality

**Migration:**
1. Create "default" project
2. Set all existing time_codes.project_id = "default"
3. User moves codes to proper projects via API/TUI
