# THE GOme/Game of life

This is a game of life written in go using [pixel](https://github.com/faiface/pixel) for graphics.
I mainly just did this so I can tell people I know go and finally put that cute Go sticker on my
laptop like a cool modern developer.

It has multiplayer.

## Usage

```bash
# run the host
go run cmd/main.go -own-port 8081 -partner-port 8080 -is-host true -own-udp-port 4321 -partner-udp-port 1234
# run the other player
go run cmd/main.go -own-port 8080 -partner-port 8081 -own-udp-port 1234 -partner-udp-port 4321
```

Now two people can click around on one board.
