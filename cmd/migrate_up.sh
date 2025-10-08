source .env
export DATABASE_URL=$(echo "$DATABASE_URL" | tr -d '\r\n')
goose -dir sql/schema postgres $DATABASE_URL up