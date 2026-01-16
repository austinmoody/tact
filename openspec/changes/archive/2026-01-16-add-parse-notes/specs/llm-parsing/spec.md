# llm-parsing Specification Delta

## ADDED Requirements

### Requirement: Parse Notes

The LLM SHALL provide reasoning notes explaining how parsing decisions were made.

#### Scenario: Successful parse with context match

- Given: Raw text "2h APHL meeting"
- And: Context document exists: "ALL meetings with APHL go to FEDS-163"
- When: The entry is parsed
- Then: `parse_notes` contains the LLM's reasoning for the match
- And: `parse_notes` includes which context document was used
- And: `parse_notes` includes the similarity score of the matched context

#### Scenario: Needs review due to ambiguity

- Given: Raw text "1h APHL UI work"
- And: Context exists: "ALL APHL meetings go to FEDS-163"
- And: Context exists: "ALL UI work goes to FEDS-167"
- When: The entry is parsed with low confidence
- Then: `parse_notes` explains the ambiguity
- And: `parse_notes` lists the conflicting context rules considered

#### Scenario: No matching context

- Given: Raw text "30m general admin"
- And: No relevant context documents match
- When: The entry is parsed
- Then: `parse_notes` explains matching was based on time code descriptions/keywords
- And: `parse_notes` indicates no specific context rules applied
