# ByteBucket

ByteBucket is a small file-storage and authentication backend written in Go. It demonstrates common backend features you'd expect in a modern web service: secure password hashing, JWT-based authentication with refresh tokens, DB-backed models using sqlc, and presigned S3 uploads for direct-to-cloud file uploads.

This repository is a work-in-progress created as a student project. It is suitable as a portfolio example to show practical skills wiring authentication, storage, and cloud integrations together. Several features are intentionally incomplete and marked under "Known limitations" below.

## Highlights / Functionality

- User registration and login
  - Passwords hashed with bcrypt (`internal/auth/hash.go`).
  - Login returns a JWT access token and a refresh token stored in the DB.
- JWT-based authentication middleware
  - Access tokens created with `github.com/golang-jwt/jwt/v5` (`internal/auth/jwt.go`) and validated in `cmd/jwt_middleware.go`.
- Refresh tokens
  - Server-generated refresh tokens stored in the DB (`internal/auth/refresh_token.go` + `sql/refresh_tokens.sql`).
- Presigned S3 uploads
  - Server generates presigned PUT URLs so clients upload directly to S3 without routing bytes through the app (`cmd/uploadFileHandler.go`).
- Database access with sqlc
  - SQL schemas in `sql/schema/` and generated accessors in `internal/database`.
- Basic API utilities
  - JSON helpers and consistent response functions in `cmd/json.go`.
- CI and tests
  - GitHub Actions workflow checks formatting, runs `go vet`, and runs `go test` with the race detector and coverage (see `.github/workflows/ci.yml`).

## Technologies Used

- Language: Go 1.24
- Authentication: bcrypt (`golang.org/x/crypto/bcrypt`), JWT (`github.com/golang-jwt/jwt/v5`)
- Database: PostgreSQL (sql files under `sql/schema`), accessed via sqlc-generated code (`internal/database`)
- Cloud Storage: Amazon S3 presigned uploads (`github.com/aws/aws-sdk-go-v2`)
- CI: GitHub Actions
- Utilities: `github.com/google/uuid`, `github.com/joho/godotenv` (for local env loading)

## Architecture Overview

1. Client registers a user (POST /users). Password is hashed and user stored in the DB.
2. Client logs in (POST /login). On success the server returns a short-lived JWT and a refresh token (longer-lived).
3. Client uses JWT to access protected endpoints. Middleware extracts the JWT, validates it, and sets the user ID in request context.
4. For file uploads, the client requests a presigned upload URL (POST /upload). Server returns a presigned PUT URL; the client uploads directly to S3 and the server records metadata in the DB.

## How to run (development)

1. Install Go 1.24 and PostgreSQL.
2. Copy `.env.example` to `.env` (if present) and set DB and AWS credentials/environment variables (e.g., `DATABASE_URL`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `S3_BUCKET`, `TOKEN_SECRET`).
3. Run DB migrations in `sql/schema/` to create users/refresh_tokens/files tables.
4. Build and run:

```powershell
go build ./cmd
.\cmd\your-binary-name.exe
```

Or run directly with `go run` from the `cmd` folder if you prefer.

API endpoints and handlers are implemented in `cmd/` (for example: `createUserHandler.go`, `loginUserHandler.go`, `uploadFileHandler.go`).

## Tests & CI

- Run unit tests locally:

```powershell
go test ./... -v
```

- CI: GitHub Actions runs formatting checks, `go vet`, and `go test` with the race detector and coverage. Coverage is uploaded as an artifact.

## Known limitations (Work in progress)

- Downloading files is not implemented. The server currently only supports creating presigned upload URLs and storing metadata.
- User account deletion is not implemented.
- File deletion (object removal & DB cleanup) is not implemented.
- Refresh tokens are stored/created but could be hardened by storing hashed refresh tokens in the DB and supporting rotation.
- Input validation needs improvement (email/password validation with clear per-field errors).
- Error responses sometimes include internal error strings — these must be sanitized before a public release.
- More unit and integration tests are desirable, including handler tests and an integration test against a test DB.

## Security notes / Suggestions

- Do not return raw error messages in responses — use machine-readable error codes and safe human messages.
- Use short lifetimes for access tokens and rotate refresh tokens on use.
- Consider storing only hashed refresh tokens server-side to reduce damage from a DB leak.
- Add rate-limiting or temporary lockout for repeated failed login attempts.
- Ensure S3 bucket CORS is correctly configured for direct browser uploads, and use HTTP-only Secure cookies for refresh tokens when serving browsers.

## Next steps / Ideas to showcase

- Implement file downloads with presigned GETs and a demo HTML page that performs the full flow.
- Add hashed refresh tokens + rotation and tests for token flows.
- Add a small React/vanilla JS demo that demonstrates signup -> login -> presigned upload -> display uploaded file list.
- Add integration tests that run migrations on a test Postgres container and run key end-to-end flows.

## Contact / Attribution

This project was created by the author as a student portfolio project. If you want help improving any of the TODO items above, I can implement them and add tests/CI changes.
# ByteBucket
