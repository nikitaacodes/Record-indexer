# Record Indexer 
a small backend service that stores immutable records on disk and exposes them via HTTP with implemented order and integrity check.

What it does?
The service accepts records and appends them to a disk file in JSON format. 
- Records are immutable once written
- the service keeps order of record noted
- integrity is validated using SHA-256 hashing algo
- the service restarts without losing data and shuts down without failures

API Endpoints 
GET /health   ( system health check )
GET /records  (list all verfied records)
GET /record?id=<int> (fetch a single verified record by ID)
all responses are returned in JSON 

Running service locally 
From project root: 
go mod tidy 
go run ./cmd 

# storage 
write records to disk and read them back safely 