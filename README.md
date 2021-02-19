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
