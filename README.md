# ByteBucket

ByteBucket is a small file-storage and authentication backend written in Go. It demonstrates common backend features you'd expect in a modern web service: secure password hashing, JWT-based authentication with refresh tokens, DB-backed models using sqlc, and S3-backed file storage.

This README reflects the current implementation: uploads go through the server (multipart) instead of presigned URLs.

## What changed
- **Uploads**: the server now accepts multipart/form-data POST to `/api/files`. The server saves the uploaded file to a temporary file (OS temp dir), uploads it to S3 via the server's AWS credentials (PutObject), records metadata in the database, and returns the object key. No presigned PUT URLs are returned.
- **Frontend**: the demo frontend now sends files as `FormData` to `/api/files` (with `Authorization: Bearer <JWT>`) and handles downloads by fetching or opening the presigned GET URL the server generates for downloads.

## Highlights / Functionality

- **User registration and login**
  - Passwords hashed with bcrypt (`internal/auth/hash.go`).
- **JWT-based authentication middleware**
  - Access tokens created and validated with `github.com/golang-jwt/jwt/v5` (`internal/auth/jwt.go`).
- **Refresh tokens**
  - Server-generated refresh tokens stored in the DB.
- **Server-side S3 uploads**
  - Client uploads files to the server via multipart POST `/api/files` (protected). Server saves to a temp file and calls S3 PutObject. The DB is updated with object metadata.
- **List files**
  - GET `/api/files` returns a JSON array of file metadata for the authenticated user.
- **Download**
  - GET `/api/files/{id}` returns a presigned GET URL; the frontend fetches as a blob and falls back to opening the URL in a new tab if browser CORS prevents reading the response.

## Technologies Used

- Language: Go 1.24
- Authentication: bcrypt, JWT
- Database: PostgreSQL with sqlc-generated code
- Cloud Storage: Amazon S3 (`github.com/aws/aws-sdk-go-v2`)
- CI: GitHub Actions
- Utilities: `github.com/google/uuid`, `github.com/joho/godotenv`

## How to run (development)

1. Install Go 1.24 and PostgreSQL.
2. Create a `.env` file with:
```
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
TOKEN_SECRET=your_jwt_secret
AWS_REGION=eu-central-1
S3_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
PORT=8080
```

3. Run database migrations in `sql/schema/` to create the tables.
4. Build and run:
```bash
go build ./cmd
./cmd
```

The server serves a demo frontend at `/` which you can use to register, login, and upload files.

## API Endpoints

- **POST /api/users** — Register a user
  - Body: `{"username":"alice","email":"alice@example.com","password":"secret"}`
  - Response: 201 with user object

- **POST /api/login** — Login
  - Body: `{"email":"alice@example.com","password":"secret"}`
  - Response: 200 with `jwt_token` and `refresh_token`

- **POST /api/files** — Upload file (protected)
  - Body: multipart/form-data with `file` field
  - Response: 200 with file metadata and object key

- **GET /api/files** — List files (protected)
  - Response: 200 with array of file metadata

- **GET /api/files/{id}** — Get download URL (protected)
  - Response: 200 with presigned GET URL (owner only)

- **POST /api/refresh** — Refresh token
  - Body: `{"refresh_token":"..."}`
  - Response: 200 with new JWT access token

## Notes & Security

- The server must have `s3:PutObject` and `s3:GetObject` permissions on the configured bucket.
- Temporary files are removed after upload but are stored briefly in the OS temp dir; consider scanning/size limiting in production.
- Use HTTPS and secure cookie patterns if you move tokens into cookies for browser flows.
- Do not commit `.env` containing secrets to source control.

## Tests & CI

Run unit tests locally:
```bash
go test ./... -v
```

GitHub Actions runs formatting checks, `go vet`, and tests with the race detector and coverage.

## Known Limitations

- User account deletion is not implemented.
- File deletion (object removal & DB cleanup) is not implemented.
- Refresh tokens could be hardened by storing hashed tokens and supporting rotation.
- Input validation needs improvement with per-field error messages.
- Error responses should sanitize internal error strings before public release.

## Next Steps

- Add hashed refresh tokens + rotation with tests
- Implement more comprehensive unit and integration tests
- Add rate-limiting for login attempts
- Improve input validation with detailed error messages