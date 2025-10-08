source .env
goose -dir sql/schema postgres $DATABASE_URL up