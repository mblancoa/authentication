# Authentication
## Packages
```
Authentication
|-- errors
|   |-- error.go
|
|-- core
|   |-- domain.go
|   |-- service.go
|   |-- port.go
|
|-- adapter
|   |-- mongodb.go
|   |-- mongodbrepository.go
|
|-- config
|   |-- domain.go
|   |-- core.go
|   |-- mongodb.go
|
|-- main.go
```
- **core**: contains the application domain entities, the services and interfaces which implement and define the business logic.
- **adapter**: implementation of interfaces in port.go

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
