# base URL
BASE="http://localhost:8080"

# 1) Register a user
curl -s -X POST "$BASE/api/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password"}' | jq .

# 2) Login (store response)
LOGIN_R=$(curl -s -X POST "$BASE/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password"}')

# show login response
echo "$LOGIN_R" | jq .

# 3) extract access & refresh tokens (works with different field name variants)
ACCESS=$(echo "$LOGIN_R" | jq -r '.jwt_token // .JWTtoken // .jwtToken // empty')
REFRESH=$(echo "$LOGIN_R" | jq -r '.refresh_token // .RefreshToken // .refreshToken // empty')

echo "ACCESS=$ACCESS"
echo "REFRESH=$REFRESH"

# fail early if access missing
if [ -z "$ACCESS" ]; then
  echo "No access token found in login response. Aborting."
  exit 1
fi

# 4) Request a presigned PUT URL for the file you want to upload
# Replace ./test.txt with your local file path, and set mime type appropriately
FILENAME="sqlc_copy.yaml"
MIMETYPE="text/plain"

PRESIGN_R=$(curl -s -X POST "$BASE/api/files" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS" \
  -d "{\"file_name\":\"$FILENAME\",\"mime_type\":\"$MIMETYPE\"}")

# show presign response
echo "$PRESIGN_R" | jq .

# extract upload URL, object key
UPLOAD_URL=$(echo "$PRESIGN_R" | jq -r '.upload_url // .uploadUrl // empty')
KEY=$(echo "$PRESIGN_R" | jq -r '.key // empty')

echo "UPLOAD_URL=$UPLOAD_URL"
echo "KEY=$KEY"

# 5) Upload file bytes directly to the presigned URL (no Authorization header)
# Use --upload-file to send the file contents; set Content-Type to the mime type returned/passed.
# Replace ./test.txt with the path to the file you want to upload.
curl -i -X PUT \
  -H "Content-Type: $MIMETYPE" \
  --upload-file ./sqlc_copy.yaml \
  "$UPLOAD_URL"

# 6) List files for the authenticated user
LIST_R=$(curl -s -H "Authorization: Bearer $ACCESS" "$BASE/api/files")
echo "$LIST_R" | jq .

# pick the first file id (adjust index if needed)
FILE_ID=$(echo "$LIST_R" | jq -r '.[0].id // empty')
echo "FILE_ID=$FILE_ID"
if [ -z "$FILE_ID" ]; then
  echo "No file ID found in list; aborting."
  exit 1
fi

# 7) Request presigned download URL for that file id
DOWNLOAD_R=$(curl -s -H "Authorization: Bearer $ACCESS" "$BASE/api/files/$FILE_ID")
echo "$DOWNLOAD_R" | jq .

DOWNLOAD_URL=$(echo "$DOWNLOAD_R" | jq -r '.url // empty')
echo "DOWNLOAD_URL=$DOWNLOAD_URL"

# 8) Download the file using the presigned URL
# Use -L to follow redirects if any. Save to downloaded_test.txt (or any name you prefer).
curl -L -o downloaded_test.txt "$DOWNLOAD_URL"
echo "Saved downloaded_test.txt"