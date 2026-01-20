## MODIFIED Requirements

### Requirement: Background Worker

The system SHALL process pending entries in the background using a non-blocking pattern that minimizes database lock duration.

#### Scenario: Automatic processing

- Given: An entry is created with status "pending"
- When: The background worker runs (every TACT_PARSER_INTERVAL seconds)
- Then: The entry is parsed
- And: status is updated to "parsed" or "failed"
- And: parsed_at is set to the current timestamp

#### Scenario: Worker disabled

- Given: `TACT_DISABLE_WORKER=true` is configured
- When: The application starts
- Then: The background worker does not start

#### Scenario: Non-blocking parse operation

- Given: An entry is pending
- When: The background worker processes the entry
- Then: The worker fetches entry data and closes the database session
- And: The LLM is called without holding a database connection
- And: A new session is opened only to write results
- And: Other API requests are not blocked during LLM processing

#### Scenario: Entry modified during parse

- Given: An entry is pending with status "pending"
- And: The worker has fetched the entry and is calling the LLM
- When: A user modifies the entry (status changes from "pending")
- And: The LLM call completes
- Then: The worker detects the status change
- And: The parse results are discarded
- And: The user's modification is preserved

#### Scenario: Entry deleted during parse

- Given: An entry is pending
- And: The worker has fetched the entry and is calling the LLM
- When: A user deletes the entry
- And: The LLM call completes
- Then: The worker detects the entry no longer exists
- And: The parse results are discarded silently
- And: No error is raised
