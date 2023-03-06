# TODO Go API
Simple Todo list Go server. Uses Go and PostgreSQL

## API documentation
See `openapi.yaml`.

## Local setup
### Golang linters
Install [golangci-lint](https://github.com/golangci/golangci-lint):
```bash
brew install golangci-lint
```

### PostgreSQL
You will need psql binary. On MacOS install it via homebrew:
```bash
brew install libpq
```

## Development
### Local build
```bash
# install build and development dependencies
make prepare
# build app
make
```

### Local run
The following command will create and run detached Docker containers running a PostgreSQL database with an [Adminer](https://www.adminer.org/) instance
```bash
# start containers
docker compose up -d
# setup database
make reset-db
# start server
make run-server
# shutdown
docker compose down
```

## Database

### Initial setup
When the local PostgreSQL database is empty or brand new, run the following script to create the schema:

```bash
docker-compose up -d
make reset-db
```

### Create a new migration
```bash
tern -m database/iban/migrations new <name>
```

### Migrate to latest version
```bash
make migrate-up
```

### Testing database migrations locally
```bash
# migrate one version up
tern migrate -c database/migrations/tern.conf -m database/migrations -d +1

# migrate one version down
tern migrate -c database/migrations/tern.conf -m database/migrations -d -1

# migrate one version up again
tern migrate -c database/migrations/tern.conf -m database/migrations -d +1
```

If no error occurred, then the migration is fine.

## Testing
### Unit tests
No setup required.

```bash
make unit
```

### Integration tests
Requires a full local setup.

```bash
docker-compose up -d
make integration
```
