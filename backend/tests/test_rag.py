import pytest
from sqlalchemy import create_engine, event
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from tact.db.base import Base
from tact.db.models import ContextDocument, Project, TimeCode
from tact.rag.embeddings import embed_text, embedding_to_array
from tact.rag.retrieval import retrieve_similar_contexts


@pytest.fixture
def db_session():
    """Create a test database session."""
    engine = create_engine(
        "sqlite:///:memory:",
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )

    @event.listens_for(engine, "connect")
    def set_sqlite_pragma(dbapi_conn, connection_record):
        cursor = dbapi_conn.cursor()
        cursor.execute("PRAGMA foreign_keys=ON")
        cursor.close()

    Base.metadata.create_all(bind=engine)
    TestSession = sessionmaker(autocommit=False, autoflush=False, bind=engine)
    session = TestSession()

    yield session

    session.close()
    engine.dispose()


class TestEmbeddings:
    """Test embedding functions."""

    def test_embed_text_returns_bytes(self):
        embedding = embed_text("hello world")
        assert isinstance(embedding, bytes)
        assert len(embedding) > 0

    def test_embed_text_consistent(self):
        text = "test embedding consistency"
        embedding1 = embed_text(text)
        embedding2 = embed_text(text)
        assert embedding1 == embedding2

    def test_embed_text_different_for_different_texts(self):
        embedding1 = embed_text("hello world")
        embedding2 = embed_text("goodbye moon")
        assert embedding1 != embedding2

    def test_embedding_to_array_roundtrip(self):
        import numpy as np

        text = "test roundtrip"
        embedding_bytes = embed_text(text)
        array = embedding_to_array(embedding_bytes)

        assert isinstance(array, np.ndarray)
        assert array.dtype == np.float32
        assert len(array) == 384  # all-MiniLM-L6-v2 dimension

    def test_embedding_is_normalized(self):
        import numpy as np

        text = "normalized embedding test"
        embedding_bytes = embed_text(text)
        array = embedding_to_array(embedding_bytes)

        norm = np.linalg.norm(array)
        assert abs(norm - 1.0) < 0.001  # Should be approximately unit length


class TestRetrieval:
    """Test context retrieval."""

    def test_retrieve_no_contexts(self, db_session):
        results = retrieve_similar_contexts("test query", db_session)
        assert results == []

    def test_retrieve_no_embeddings(self, db_session):
        # Create project and context without embedding
        project = Project(id="test", name="Test Project")
        db_session.add(project)
        db_session.commit()

        context = ContextDocument(
            project_id="test",
            content="Test content",
            embedding=None,  # No embedding
        )
        db_session.add(context)
        db_session.commit()

        results = retrieve_similar_contexts("test query", db_session)
        assert results == []

    def test_retrieve_with_embeddings(self, db_session):
        # Create project and context with embedding
        project = Project(id="test", name="Test Project")
        db_session.add(project)
        db_session.commit()

        content = "ALL meetings with APHL go to FEDS-163"
        context = ContextDocument(
            project_id="test",
            content=content,
            embedding=embed_text(content),
        )
        db_session.add(context)
        db_session.commit()

        # Query with similar text
        results = retrieve_similar_contexts("APHL meeting", db_session)
        assert len(results) == 1
        assert results[0].content == content
        assert results[0].project_id == "test"
        assert results[0].similarity > 0.3

    def test_retrieve_similarity_ranking(self, db_session):
        # Create project and multiple contexts
        project = Project(id="test", name="Test Project")
        db_session.add(project)
        db_session.commit()

        # Add contexts with different relevance
        contexts = [
            "UI development with React",
            "Backend Python API work",
            "React component styling",
        ]
        for content in contexts:
            ctx = ContextDocument(
                project_id="test",
                content=content,
                embedding=embed_text(content),
            )
            db_session.add(ctx)
        db_session.commit()

        # Query for React-related content
        results = retrieve_similar_contexts("React frontend work", db_session, top_k=3)

        # Should return results sorted by relevance
        assert len(results) > 0
        # First result should be most similar to React
        assert "React" in results[0].content or "UI" in results[0].content

    def test_retrieve_respects_top_k(self, db_session):
        project = Project(id="test", name="Test Project")
        db_session.add(project)
        db_session.commit()

        # Add many contexts
        for i in range(10):
            content = f"Context document number {i}"
            ctx = ContextDocument(
                project_id="test",
                content=content,
                embedding=embed_text(content),
            )
            db_session.add(ctx)
        db_session.commit()

        results = retrieve_similar_contexts("document", db_session, top_k=3)
        assert len(results) <= 3

    def test_retrieve_respects_min_similarity(self, db_session):
        project = Project(id="test", name="Test Project")
        db_session.add(project)
        db_session.commit()

        # Add a context
        content = "Very specific technical documentation about React hooks"
        ctx = ContextDocument(
            project_id="test",
            content=content,
            embedding=embed_text(content),
        )
        db_session.add(ctx)
        db_session.commit()

        # Query with completely unrelated text should get low similarity
        results = retrieve_similar_contexts(
            "banana fruit salad recipe",
            db_session,
            min_similarity=0.9,  # Very high threshold
        )
        assert len(results) == 0

    def test_retrieve_includes_time_code_context(self, db_session):
        # Create project and time code
        project = Project(id="test", name="Test Project")
        db_session.add(project)
        db_session.commit()

        time_code = TimeCode(
            id="FEDS-163",
            project_id="test",
            name="Development",
            description="Dev work",
        )
        db_session.add(time_code)
        db_session.commit()

        # Add context to time code
        content = "ALL deployments go to this code"
        ctx = ContextDocument(
            time_code_id="FEDS-163",
            content=content,
            embedding=embed_text(content),
        )
        db_session.add(ctx)
        db_session.commit()

        results = retrieve_similar_contexts("deployment work", db_session)
        assert len(results) == 1
        assert results[0].time_code_id == "FEDS-163"
        assert results[0].project_id is None
