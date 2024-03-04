gen:
	sqlc generate

setup:
	createdb quinn-lifts
	psql -d quinn-lifts -a -f internal/db/schema.sql

resetdb:
	dropdb -f quinn-lifts
	createdb quinn-lifts
	psql -d quinn-lifts -a -f internal/db/schema.sql

dist:
	go build -o /usr/local/bin/quinn-lifts cmd/main.go
	mkdir -p ~/Library/LaunchAgents
	cp daemon/com.lifts.plist ~/Library/LaunchAgents

launch-agent:
	launchctl unload ~/Library/LaunchAgents/com.lifts.plist
	launchctl load ~/Library/LaunchAgents/com.lifts.plist
	launchctl start com.quinn-lifts
