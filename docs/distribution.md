you want devs to easily run the kitty-beads web UI locally against their own .beads data.

The most seamless approach: embed the server into bd as a subcommand — bd ui or bd serve. They already have bd installed, it just
opens the dashboard locally pointing at their project's .beads. No second binary to distribute.

bd ui # starts server, opens browser at localhost:8080

The alternative is a standalone kitty-beads binary they run in their project dir, but that's two things to install and keep in
sync.

Want to implement bd ui that spins up the server? The server code is already written — it'd mostly be wiring cmd/server/main.go
logic into a new bd command.
