# StreamRec Core API

Lightweight RESTful backend service for a streaming recommendation platform built with Go, Chi router, and MongoDB.

## Features
- **Clean Architecture**: Separated domain, service, repository, and transport layers.
- **NoSQL Data Layer**: MongoDB collections for users, items, and user interaction tracking.
- **Fast Queries**: Compound and unique indexes for fast lookups and sorted history retrieval.
- **Catalog Search**: Substring and category search with case-insensitive filters.
- **Interactive Documentation**: Swagger UI for endpoint testing.

---

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.22+

### 1. Run MongoDB
Start the MongoDB container:
```bash
docker compose up -d
