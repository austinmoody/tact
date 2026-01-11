## 1. Schemas

- [x] 1.1 Create `backend/src/tact/schemas/__init__.py`
- [x] 1.2 Create `backend/src/tact/schemas/time_code.py` with create/update/response models
- [x] 1.3 Create `backend/src/tact/schemas/work_type.py` with create/update/response models

## 2. Time Codes Routes

- [x] 2.1 Create `backend/src/tact/routes/time_codes.py`
- [x] 2.2 Implement POST `/time-codes` - create time code
- [x] 2.3 Implement GET `/time-codes` - list time codes (with optional `active` filter)
- [x] 2.4 Implement GET `/time-codes/{id}` - get single time code
- [x] 2.5 Implement PUT `/time-codes/{id}` - update time code
- [x] 2.6 Implement DELETE `/time-codes/{id}` - soft-delete (set active=false)
- [x] 2.7 Register router in main.py

## 3. Work Types Routes

- [x] 3.1 Create `backend/src/tact/routes/work_types.py`
- [x] 3.2 Implement POST `/work-types` - create work type
- [x] 3.3 Implement GET `/work-types` - list work types (with optional `active` filter)
- [x] 3.4 Implement GET `/work-types/{id}` - get single work type
- [x] 3.5 Implement PUT `/work-types/{id}` - update work type
- [x] 3.6 Implement DELETE `/work-types/{id}` - soft-delete (set active=false)
- [x] 3.7 Register router in main.py

## 4. Testing

- [x] 4.1 Create `backend/tests/test_time_codes.py` with full CRUD tests
- [x] 4.2 Create `backend/tests/test_work_types.py` with full CRUD tests
- [x] 4.3 Verify `make test` passes
- [x] 4.4 Verify `make lint` passes

## Verification

Steps to independently verify the implementation:

1. **Start the app:**
   ```bash
   make run
   ```

2. **Test Time Codes CRUD:**
   ```bash
   # Create
   curl -X POST http://localhost:2100/time-codes \
     -H "Content-Type: application/json" \
     -d '{"id": "PROJ-001", "name": "Project Alpha", "description": "Main project"}'

   # List
   curl http://localhost:2100/time-codes

   # Get single
   curl http://localhost:2100/time-codes/PROJ-001

   # Update
   curl -X PUT http://localhost:2100/time-codes/PROJ-001 \
     -H "Content-Type: application/json" \
     -d '{"name": "Project Alpha Updated"}'

   # Delete (soft)
   curl -X DELETE http://localhost:2100/time-codes/PROJ-001

   # Verify inactive
   curl http://localhost:2100/time-codes/PROJ-001
   # Should show active=false
   ```

3. **Test Work Types CRUD:**
   ```bash
   # Create
   curl -X POST http://localhost:2100/work-types \
     -H "Content-Type: application/json" \
     -d '{"id": "dev", "name": "Development"}'

   # List
   curl http://localhost:2100/work-types

   # Get/Update/Delete follow same pattern as time-codes
   ```

4. **Run tests:**
   ```bash
   make test
   ```
