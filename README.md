# Authentication
This project is intended to be a set of proofs of concept, so its business logic should not be considered
## Packages
```
Authentication
|-- errors
|   |-- error.go
|
|-- core
|   |-- ports
|       |-- notifications.go
|       |-- persistence.go
|       |-- cache.go
|   |-- domain
|       |-- domain.go
|   |-- config.go
|   |-- service.go
|
|-- adapters
|   |--mongodb
|      |-- config.go
|      |-- service.go //persistence service implementation
|      |-- repository.go
|   |--cache
|      |-- config.go
|      |-- client.go 
|   |--notifications
|      |-- config.go
|      |-- service.go 
|
|-- api
|   |-- controllers
|       |-- base.go
|       |-- auth.go
|       |-- domain.go
|   |-- config.go
|   |-- router.go
|
|-- main.go
```

- **core**: contains the application domain entities, the services and interfaces which implement and define the business logic.
- **adapters**: implementation of interfaces being in port.go
- **api**: API implementation

## Repositories generation

Installation 
>`go install github.com/sunboyy/repogen@latest`

Generation
>`make code-generation`

## Mocks generation
Installation
>`go install github.com/vektra/mockery/v2@v2.40.1`

Generation
>`make clean mocks`
