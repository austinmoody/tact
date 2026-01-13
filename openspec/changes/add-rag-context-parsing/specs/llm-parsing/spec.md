# llm-parsing Specification Delta

## MODIFIED Requirements

### Requirement: Parse Entry Fields

The LLM SHALL extract structured fields from raw entry text using RAG-retrieved context.

#### Scenario: Context-aware matching

- Given: Raw text "2h APHL meeting about UI"
- And: Context document exists: "ALL meetings with APHL go to FEDS-163 regardless of topic"
- And: Context document exists: "ALL UI work goes to FEDS-167"
- When: The entry is parsed
- Then: Relevant context is retrieved via vector similarity
- And: time_code_id is set to "FEDS-163" (APHL rule overrides UI)
- And: The LLM uses context to make the disambiguation

#### Scenario: Acronym expansion

- Given: Raw text "1h IZG deployment"
- And: Project context exists: "IZG = IZ Gateway"
- And: Time code context exists for FEDS-163: "ALL deployments go to this code"
- When: The entry is parsed
- Then: The LLM understands IZG refers to IZ Gateway
- And: time_code_id is set to "FEDS-163"

#### Scenario: No relevant context

- Given: Raw text "30m general admin work"
- And: No specific context matches this entry
- When: The entry is parsed
- Then: The LLM falls back to time code descriptions and keywords
- And: A reasonable match is made based on available information

## ADDED Requirements

### Requirement: RAG Context Retrieval

The parser SHALL retrieve relevant context documents before calling the LLM.

#### Scenario: Retrieve similar context

- Given: An entry "2h security scan review"
- And: Context documents exist with various content
- When: The parser prepares the LLM prompt
- Then: The entry text is embedded using the same model as context docs
- And: Top-k most similar context chunks are retrieved
- And: Retrieved chunks are included in the LLM prompt

#### Scenario: Context includes source

- Given: Context is retrieved for an entry
- When: The LLM prompt is built
- Then: Each context chunk is labeled with its source (project or time code)
- And: The LLM can see which time code each rule applies to

#### Scenario: Empty context store

- Given: No context documents exist in the system
- When: An entry is parsed
- Then: Parsing proceeds without RAG context
- And: The LLM uses only time code descriptions and keywords

### Requirement: Local Embeddings

The system SHALL generate embeddings locally without external API calls.

#### Scenario: Embed entry text

- Given: An entry is submitted for parsing
- When: The parser retrieves context
- Then: The entry text is embedded using a local model
- And: No external embedding API is called

#### Scenario: Embed context document

- Given: A context document is created
- When: The document is saved
- Then: The content is embedded using a local model
- And: The embedding is stored with the document
