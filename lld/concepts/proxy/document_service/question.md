
## Problem Statement

You are building a document storage service for an internal enterprise system.

The system allows users to:
- View documents
- Download documents

However:
- Documents are large (PDFs, scans, contracts)
- Fetching a document from storage is expensive
- Not all users are authorized to access all documents
- Access must be audited


## Approach

- domcument service
    - fetch the document from storage
    - return the document
- proxy document service
    - check the user based access
    - Add a log for audit - async
    - cache the document
    - calls the real document service

