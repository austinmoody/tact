## 1. Schemas

- [ ] 1.1 Create `backend/src/tact/schemas/entry.py` with Pydantic models:
  - `EntryCreate` (raw_text required, entry_date optional)
  - `EntryUpdate` (all fields optional for PATCH)
  - `EntryResponse` (full entry representation)
- [ ] 1.2 Export schemas from `backend/src/tact/schemas/__init__.py`

## 2. Entries Routes

- [ ] 2.1 Create `backend/src/tact/routes/entries.py`
- [ ] 2.2 Implement POST `/entries` - create entry (status=pending, entry_date defaults to today)
- [ ] 2.3 Implement GET `/entries` - list entries with filters:
  - `status` filter
  - `time_code_id` filter
  - `work_type_id` filter
  - `from_date` / `to_date` range filter
  - `limit` / `offset` pagination
- [ ] 2.4 Implement GET `/entries/{id}` - get single entry
- [ ] 2.5 Implement PATCH `/entries/{id}` - update entry (sets manually_corrected=true)
- [ ] 2.6 Implement DELETE `/entries/{id}` - hard delete entry (returns 204)
- [ ] 2.7 Register router in main.py

## 3. Testing

- [ ] 3.1 Create `backend/tests/test_entries.py` with tests for:
  - Create entry (raw_text only, with entry_date, missing raw_text)
  - List entries (all, with each filter, pagination)
  - Get single entry (exists, not found)
  - Update entry (success, sets manually_corrected, not found)
  - Delete entry (success, not found)
- [ ] 3.2 Verify `make test` passes
- [ ] 3.3 Verify `make lint` passes

## 4. Documentation

- [ ] 4.1 Update backend/README.md with entries endpoints

## Verification

Steps to independently verify the implementation:

1. **Start the app:**
   ```bash
   make run
   ```

2. **Test Entries CRUD:**
   ```bash
   # Create entry (minimal - just raw_text)
   curl -X POST http://localhost:2100/entries \
     -H "Content-Type: application/json" \
     -d '{"raw_text": "2h coding on Project Alpha"}'
   # Response includes generated id, status="pending", entry_date=today

   # Create entry with specific date
   curl -X POST http://localhost:2100/entries \
     -H "Content-Type: application/json" \
     -d '{"raw_text": "1h meeting yesterday", "entry_date": "2026-01-09"}'

   # List all entries
   curl http://localhost:2100/entries

   # List with filters
   curl "http://localhost:2100/entries?status=pending"
   curl "http://localhost:2100/entries?from_date=2026-01-01&to_date=2026-01-31"
   curl "http://localhost:2100/entries?limit=10&offset=0"

   # Get single entry (use id from create response)
   curl http://localhost:2100/entries/{id}

   # Update entry
   curl -X PATCH http://localhost:2100/entries/{id} \
     -H "Content-Type: application/json" \
     -d '{"duration_minutes": 120}'
   # Response shows manually_corrected=true

   # Delete entry
   curl -X DELETE http://localhost:2100/entries/{id}
   # Returns 204 No Content
   ```

3. **Run tests:**
   ```bash
   make test
   ```
