import logging
from dataclasses import dataclass

import numpy as np
from sqlalchemy.orm import Session

from tact.db.models import ContextDocument
from tact.rag.embeddings import embed_text, embedding_to_array

logger = logging.getLogger(__name__)


@dataclass
class RetrievedContext:
    """A retrieved context document with its similarity score."""

    id: str
    content: str
    project_id: str | None
    time_code_id: str | None
    similarity: float


def retrieve_similar_contexts(
    query: str,
    session: Session,
    top_k: int = 5,
    min_similarity: float = 0.3,
) -> list[RetrievedContext]:
    """
    Retrieve the most similar context documents for a query.

    Args:
        query: The text to search for similar contexts
        session: Database session
        top_k: Maximum number of results to return
        min_similarity: Minimum cosine similarity threshold

    Returns:
        List of retrieved contexts sorted by similarity (highest first)
    """
    # Embed the query
    query_embedding = embedding_to_array(embed_text(query))

    # Get all context documents with embeddings
    contexts = (
        session.query(ContextDocument)
        .filter(ContextDocument.embedding.isnot(None))
        .all()
    )

    if not contexts:
        logger.debug("No context documents with embeddings found")
        return []

    # Calculate cosine similarities
    results = []
    for ctx in contexts:
        ctx_embedding = embedding_to_array(ctx.embedding)
        # Both embeddings are normalized, so dot product = cosine similarity
        similarity = float(np.dot(query_embedding, ctx_embedding))

        if similarity >= min_similarity:
            results.append(
                RetrievedContext(
                    id=ctx.id,
                    content=ctx.content,
                    project_id=ctx.project_id,
                    time_code_id=ctx.time_code_id,
                    similarity=similarity,
                )
            )

    # Sort by similarity (highest first) and limit to top_k
    results.sort(key=lambda x: x.similarity, reverse=True)
    results = results[:top_k]

    logger.debug(
        "Retrieved %d context documents for query (top_k=%d, min_sim=%.2f)",
        len(results),
        top_k,
        min_similarity,
    )
    return results
