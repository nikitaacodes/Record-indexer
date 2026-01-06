# Record Indexer 
a small backend service that stores immutable records on disk and exposes them via HTTP with implemented order and integrity check.

What it does?
The service accepts records and appends them to a disk file in JSON format. 
- Records are immutable once written
- the service keeps order of record noted
- integrity is validated using SHA-256 hashing algo
- the service restarts without losing data and and clean shutdowns.

Security & Trust 
- env-based Basic Auth 
- In-memory per-IP rate limiting to prevent spam bursts
- /health : public route
- all other routes are protected
- Automatic cleanup of old rate-limit entries from memory 
 
API Endpoints 
GET /health    |      Public     |   ( system health check )
GET /records?limit=&offset= |  protected |(paginated list of integrity verified records)
GET /record?id=<int>   |  Protected   |  (fetch a single verified record by ID)
GET /integrity/status    | protected  |  returns integrity summary in JSON
response format:
{
    "total_records" : 5,
    "valid_records" : 2,
    "corrupted_records" : 1,
    "last_checked": "2026-01-05T..."
}



Running service locally 
From project root: 
go mod tidy 
go run ./cmd 

# storage 
write records to disk and read them back safely 