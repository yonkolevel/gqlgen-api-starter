## Intro

Based off of - https://github.com/cmelgarejo/go-gql-server

### Postgres db setup

## Introduction

Heavily inspired on https://dev.to/cmelgarejo/creating-an-opinionated-graphql-server-with-go-part-1-3g3l
https://github.com/cmelgarejo/go-gql-server/tree/master

## DB Setup

```bash
docker run -p 5433:5432 --name gqlgen-api-starter-db -e POSTGRES_PASSWORD=password -e POSTGRES_USER=admin -e POSTGRES_DB=gqlgen-api-starter -v ${HOME}/postgres-data/gqlgen-api-starter:/var/lib/postgresql/data -d postgres
```



## Running the server

### Env variables

- Rename the .env-example file to .env
- Request the secrets to a member of the team

```bash
go get

./scripts/gqlgen.sh

./scripts/build.sh

./scripts/dev-run.sh
```

Alternatively you can run these tasks with a task runner called [Task](https://taskfile.dev/#/installation):

```bash
brew install go-task/tap/go-task

// or using go

go install github.com/go-task/task/v3/cmd/task@latest
```

Please read the `Taskfile.yml` to see available commands. Some examples are:

```bash
task setup
task dev
task run
task build
```

Use `task dev --watch` for improved dev experience, where the server is reloaded as you make changes without having to stop and run again.
