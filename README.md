# Authentication
## Modules
- **core**: contains the application domain entities, the services and interfaces which implement and define the business logic.
- **cache-redis**: cache implementation using redis
- **repository-mongodb**: repository implementation to persist data into a mongodb database
- **api-server**: application which runs in a server and provides the access through an API
- **api-client**: library which provides the methods to connect with the API

## Mocks generation
Installation

>`go install github.com/vektra/mockery/v2@v2.40.1`

Generation

>`mockery`