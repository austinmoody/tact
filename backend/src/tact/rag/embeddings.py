import logging
from functools import lru_cache

import numpy as np
from sentence_transformers import SentenceTransformer

logger = logging.getLogger(__name__)

# Model name for embeddings - all-MiniLM-L6-v2 is a good balance of size and quality
MODEL_NAME = "all-MiniLM-L6-v2"
EMBEDDING_DIMENSION = 384  # Dimension for all-MiniLM-L6-v2


@lru_cache(maxsize=1)
def get_embedding_model() -> SentenceTransformer:
    """Load and cache the sentence transformer model."""
    logger.info("Loading embedding model: %s", MODEL_NAME)
    model = SentenceTransformer(MODEL_NAME)
    logger.info("Embedding model loaded successfully")
    return model


def embed_text(text: str) -> bytes:
    """
    Generate an embedding for the given text.

    Args:
        text: The text to embed

    Returns:
        The embedding as bytes (numpy array serialized)
    """
    model = get_embedding_model()
    embedding = model.encode(text, convert_to_numpy=True)
    # Normalize the embedding for cosine similarity
    embedding = embedding / np.linalg.norm(embedding)
    return embedding.astype(np.float32).tobytes()


def embedding_to_array(embedding_bytes: bytes) -> np.ndarray:
    """
    Convert embedding bytes back to numpy array.

    Args:
        embedding_bytes: The embedding as bytes

    Returns:
        The embedding as a numpy array
    """
    return np.frombuffer(embedding_bytes, dtype=np.float32)
