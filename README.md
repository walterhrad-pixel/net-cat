# NetCat-like TCP Chat Application

A fully working TCP chat server written in Go with clean architecture, supporting multiple rooms, user commands, and chat history.

## Features

- **TCP Server**: Listens on configurable port (default: 8989)
- **Multiple Clients**: Supports up to 10 clients per room
- **Multiple Rooms**: Create and join rooms with `/join` command
- **User Commands**:
  - `/name newname` - Change your username
  - `/join roomName` - Join or create a room
- **Chat History**: Stores and replays last 100 messages when joining a room
- **System Messages**: Join/leave notifications
- **Timestamped Messages**: All messages include timestamps
- **ASCII Welcome Banner**: Linux-style welcome screen

## Project Structure

```
net-cat/
├── main.go           # Entry point with flag parsing
├── server/
│   ├── server.go     # TCP server implementation
│   ├── format.go     # Message formatting utilities
│   ├── client.go     # Client connection handling
│   ├── hub.go        # Central hub for room management
│   ├── room.go       # Room logic and broadcasting
│   └── message.go    # Message types and serialization
├── utils/
│   └── ascii.go      # ASCII art welcome banner
└── go.mod
├── README.md
```

## Architecture

- **Hub**: Manages all rooms and coordinates client movement between rooms
- **Room**: Manages clients in a room, broadcasts messages, stores history
- **Client**: Handles individual connection, reads/writes, processes commands
- **Message**: Structured message with type, username, content, room, timestamp

Communication flow:
1. Server accepts TCP connections
2. Client provides username
3. Client joins default "general" room
4. Each client runs in two goroutines (read/write)
5. Messages flow through room broadcast channels
6. Hub coordinates room creation and client transfers

## Usage

### Start the Server

```bash
go run .
```

Or with custom flags:

```bash
go run . -p 9999 -room lobby
```

### Connect with NetCat

```bash
nc localhost 8989
```

You'll be prompted for a username. After entering your name, you can:

- Send chat messages (just type and press Enter)
- Change name: `/name NewUsername`
- Join room: `/join roomName`
- Create room: `/join newroom` (creates if doesn't exist)

### Example Session

```
Enter your username: Alice
[2026-05-04 10:00:00][SYSTEM]: Welcome Alice! You joined room 'general'.
[2026-05-04 10:00:00][SYSTEM]: Alice has joined the room
--- Chat History ---
--- End of History ---

Hello everyone!
[2026-05-04 10:00:05][SYSTEM]: Bob has joined the room
Hi Bob!
/name AliceSmith
[2026-05-04 10:00:10][SYSTEM]: Alice changed name to AliceSmith
/join tech
[2026-05-04 10:00:15][SYSTEM]: AliceSmith has left the room
[2026-05-04 10:00:15][SYSTEM]: AliceSmith joined room 'tech'
```

## Running Tests

The application can be tested with multiple netcat clients:

Terminal 1:
```bash
./netcat -p 8989
```

Terminal 2:
```bash
nc localhost 8989
# Enter: Alice
# Type: Hello!
```

Terminal 3:
```bash
nc localhost 8989
# Enter: Bob
# Type: Hi Alice!
```

Both clients will see each other's messages.

## Build

```bash
go build -o netcat .
```

Then run:

```bash
./netcat
```

## Code Quality

- Small, focused functions (<20 lines)
- No duplicated logic
- Clean naming conventions
- Proper error handling
- No deadlocks (uses channels and mutexes correctly)
- Goroutines for concurrent client handling

## Technical Details

- **Concurrency**: Each client has dedicated read/write goroutines
- **Synchronization**: Mutexes protect shared state in rooms
- **Channels**: Message passing between components
- **History**: Circular buffer of last 100 messages per room
- **No Echo**: Sender doesn't receive their own messages
- **Empty Message Filtering**: Ignores empty input

## Future Enhancements

- Terminal UI with gocui
- Private messaging
- Room listing
- Kick/ban functionality
- Message persistence
- TLS/SSL support