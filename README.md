# simple-chat-server
Simple Socket Based Chat Server

## Starting Server
Start Server using command `go run chatserver.go`

## Starting Client
Use `nc localhost 4000` via terminal to initiate session with Server (as a client)

## Current Features
- Matches first available clients with each other.
- If one client disconnects, puts the other client back into the matching pool.
- Allows multiple simultaneous chat sessions.
