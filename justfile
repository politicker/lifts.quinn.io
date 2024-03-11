gen:
	sqlc generate

setup:
	createdb quinn-lifts
	psql -d quinn-lifts -a -f internal/db/schema.sql

resetdb:
	dropdb -f quinn-lifts
	createdb quinn-lifts
	psql -d quinn-lifts -a -f internal/db/schema.sql

dist-mac:
	go build -o /usr/local/bin/quinn-lifts cmd/main.go
	mkdir -p ~/Library/LaunchAgents
	cp daemon/com.lifts.plist ~/Library/LaunchAgents

dist-linux:
	go build -o /home/quinn/.config/bin/quinn-lifts cmd/main.go
	cp daemon/lifts.service /home/quinn/.config/systemd/user/lifts.service

systemd:
	systemctl --user enable lifts
	systemctl --user start lifts

systemd-reload:
	systemctl --user daemon-reload
	systemctl --user restart lifts

check:
	systemctl --user status lifts

launch-agent:
	launchctl unload ~/Library/LaunchAgents/com.lifts.plist
	launchctl load ~/Library/LaunchAgents/com.lifts.plist
	launchctl start com.quinn-lifts
