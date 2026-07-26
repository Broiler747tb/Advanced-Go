# Running the Order API

The app needs a PostgreSQL database. If you don't have Postgres installed
locally, a **free hosted database** is the fastest way to run it — no install,
and it works even when your local disk is full.

## 1. Get a free Postgres (2 minutes)

1. Go to <https://neon.tech> (or <https://supabase.com>) and sign up (free tier).
2. Create a project / database.
3. Copy the **connection string**. It looks like:

   ```
   postgresql://myuser:mypassword@ep-cool-name-123.eu-central-1.aws.neon.tech/neondb?sslmode=require
   ```

## 2. Point the app at it

Open `.env` and replace the `DSN` line with your connection string:

```env
DSN="postgresql://myuser:mypassword@ep-cool-name-123.../neondb?sslmode=require"
TOKEN="dev-only-change-me-9f3a1c7e5b2d4680aa11bb22cc33dd44"
```

GORM accepts both the URL form above and the `host=... user=... password=...`
key/value form — either works. `TOKEN` is the secret used to sign JWTs.

## 3. Run it

The server auto-creates its tables on startup, so a single command is enough:

```bash
go run ./Cmd
# Listening and serving on port: 8081
```

## 4. Try the SMS auth flow

```bash
# 1. Request a code — returns {"sessionId":"..."}
curl -s -X POST localhost:8081/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone":"89990009900"}'

# The code is printed in the server console, e.g.:
#   [SMS] phone=89990009900 code=4500 session=21b8...

# 2. Verify the code — returns {"token":"<jwt>"}
curl -s -X POST localhost:8081/auth/verify-code \
  -H "Content-Type: application/json" \
  -d '{"sessionId":"<sessionId from step 1>","code":4500}'

# 3. Call a protected endpoint with the token
curl -s -X POST localhost:8081/link \
  -H "Authorization: Bearer <token from step 2>" \
  -H "Content-Type: application/json" \
  -d '{"name":"my link"}'
```

## Running the tests (no database needed)

The full auth flow is covered by tests that use an in-memory fake, so they run
anywhere with no setup:

```bash
go test ./...
```
