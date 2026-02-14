# How To

## Env

Rename example.env to .env and change the value based on requierment.

## DB Migration

DB name tlab_wallet

`export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/tlab_wallet?sslmode=disable'
`
`migrate -database ${POSTGRESQL_URL} -path sql/migrations up`

## Run

in root folder aka inside tlab-wallet foler run:
`go run ./cmd/`

## Doc

All endpoint requires Bearer token but these two:

```
POST /api/auth/register
```

and

```
POST /api/auth/login
```

For more endpoint [Check this out!](https://github.com/rianprayoga/tlab-wallet/blob/main/tlab-documentation.html)
