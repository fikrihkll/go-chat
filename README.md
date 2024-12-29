# Golang Chat Service API
This application uses Echo framework for HTTP request. This project has unit test to verify business logics, and for easier integration testing this application uses dockertest (will be implemented soon).

The Chat Application (Fluterr App) will mainly use MQTT to listen live update Message from the chat room. And the API will take part to save the message history.

### Why use MQTT?
The actual purpose of using MQTT to handle live update message is just to adapt with MQTT technology.

### Why use Golang as the backend?
This is my very first Golang project and this project is aimed for the preparation before joining Pinhome :D. Basically, I want to learn and get used to Golang syntax, commands and its architectures.

## Step To run the API

- copy `env.example` to `.env`
- run 
```bash
make docker
```
- run 
```bash
make migrate
```
- run 
```bash
make run
```

## To run integration test use this command :

```bash
make test
```


## Makefile scripts
this project conatains Makefile :

|command|description|
|---|---|
| make migrate | migrate up database to writedb |
| make migrate-down | migrate down database from writedb |
| make migrate-step-down | migrate down database 1 version from writedb |
| make migrate-info | shows migration info such as dirty status and version |
| make migrate-force version={version_to_force} | force migration to a specific version |
| make seed | seed/insert sample data to system |
| make swagger | rebuild swagger api docs |
| make docker | compose up minimum dependencies (all single node) |
| make run | run api |



## API Docs
Open swagger ui at `http://localhost:3000/swagger/index.html`

