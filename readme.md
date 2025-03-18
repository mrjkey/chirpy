## install air

```bash
 go install github.com/air-verse/air@latest
 air
```

this is just helpful when working on a server, it rebuilds whenever you change code

- add tmp/ to gitignore

if you want to change any of the default behavior, you can configure the air.toml file in the top level dir

example

```toml
[build]
cmd = "go build -o tmp/main ."
bin = "tmp/main"
full_bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["vendor", "tmp"]
[log]
time = false
```

for now i will just leave the defaults and see how it goes. what happens if I change the readme?
nah, doesn't reload, good

## using extension

it is just called rest client

if you can get postman to work, use that

or curl

## connection string

```bash
psql "postgres://postgres:a@localhost:5432/chirpy"
```

migrations

```bash
cd sql/schema
goose postgres "postgres://postgres:a@localhost:5432/chirpy" up
goose postgres "postgres://postgres:a@localhost:5432/chirpy" down
```

## .env file

it is set to ignore (because you should it may have secrets)
but the content is this:

DB_URL="database_url" (the postgres url)
PLATFORM="dev"
SECRET_TOKEN="someshit"

## stuff to install

go
sqlc
goose
postgres
gotdotenv
air

go get github.com/joho/godotenv

## making a secret

```bash
openssl rand -base64 64
```
