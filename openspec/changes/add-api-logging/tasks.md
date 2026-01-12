# Tasks: Add Structured API Logging

## 1. Logging Configuration

- [ ] 1.1 Configure structured log format with timestamp, level, module, and message
- [ ] 1.2 Add correlation ID generation for request tracing
- [ ] 1.3 Ensure logs output to stdout for Docker/container compatibility

## 2. Request Middleware

- [ ] 2.1 Create request logging middleware to capture method, path, status, and timing
- [ ] 2.2 Attach correlation ID to each request context
- [ ] 2.3 Register middleware in FastAPI application

## 3. Entries Route Logging

- [ ] 3.1 Add logging to `create_entry` - log entry creation with ID
- [ ] 3.2 Add logging to `list_entries` - log query parameters and result count
- [ ] 3.3 Add logging to `get_entry` - log entry retrieval
- [ ] 3.4 Add logging to `update_entry` - log entry modification with changed fields
- [ ] 3.5 Add logging to `delete_entry` - log entry deletion
- [ ] 3.6 Add logging to `reparse_entry` - log reparse request

## 4. Work Types Route Logging

- [ ] 4.1 Add logging to work type CRUD operations
- [ ] 4.2 Log create, update, and delete operations with relevant identifiers

## 5. Time Codes Route Logging

- [ ] 5.1 Add logging to time code CRUD operations
- [ ] 5.2 Log create, update, and delete operations with relevant identifiers

## 6. Testing

- [ ] 6.1 Verify logs appear in Docker output when running API
- [ ] 6.2 Confirm log entries contain expected information for each operation

## Verification

Steps to verify the implementation works:

1. Start the API with `docker compose up`
2. Create an entry via POST `/entries` and verify log output shows the creation
3. List entries via GET `/entries` and verify log shows query and count
4. Update an entry via PATCH `/entries/{id}` and verify log shows modification
5. Delete an entry via DELETE `/entries/{id}` and verify log shows deletion
6. Create/update/delete work types and time codes, verifying each is logged
7. Check that all logs include timestamps and consistent formatting
