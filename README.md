# How To

## Env

Rename example.env to .env and change the value based on requierment.

## DB Migration

DB name tlab_wallet

`export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/tlab_wallet?sslmode=disable'
`
`migrate -database ${POSTGRESQL_URL} -path sql/migrations up`

## Run

in root foler aka inside tlab-wallet foler run:
`go run ./cmd/`

## Doc

All endpoint but requiers Bearer token but this one:

```
POST /api/auth/register
```

and

```
POST /api/auth/login
```

[Check this out!](https://github.com/rianprayoga/tlab-wallet/blob/main/tlab-documentation.html)
