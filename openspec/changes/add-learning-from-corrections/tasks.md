# Tasks

## 1. Update Entry Update Endpoint

- [x] 1.1 Add `learn: bool = True` query parameter to PATCH `/entries/{id}`
- [x] 1.2 After successful update, check if `learn=True` and entry has `time_code_id`
- [x] 1.3 If learning, create context document for the time code with formatted content
- [x] 1.4 Log when a learned context document is created

## 2. Add Tests

- [x] 2.1 Test that updating entry with `learn=true` creates context document
- [x] 2.2 Test that updating entry with `learn=false` does not create context
- [x] 2.3 Test that updating entry without time_code_id does not create context
- [x] 2.4 Test that context content format is correct

## 3. Documentation

- [x] 3.1 Update backend README API documentation for the new parameter
- [x] 3.2 Add example curl command showing learning behavior

## Dependencies

- Task 2 depends on Task 1 (need implementation to test)
- Task 3 can be done in parallel with Tasks 1 and 2

## Notes

- The embedding will be automatically generated when creating the context document (existing behavior)
- Context document content format:
  ```
  Example: "{raw_text}"
  Parsed as: {duration_minutes} minutes, work_type: {work_type_id}
  ```
- If duration_minutes is None, omit that part from the content
- If work_type_id is None, omit that part from the content
