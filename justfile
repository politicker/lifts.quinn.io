gen:
	sqlc generate

resetdb:
	dropdb -f quinn-lifts
	createdb quinn-lifts
	psql -d quinn-lifts -a -f internal/db/schema.sql
