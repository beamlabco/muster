# Muster CLI

> A beautiful terminal user interface (TUI) for Muster - a team productivity tool for daily standups, attendance tracking, and leave management.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Node.js](https://img.shields.io/badge/Node.js-18+-green.svg)](https://nodejs.org/)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)

---

## Features

- 🔐 **Authentication**: Login and registration with beautiful forms
- 📝 **Standups**: Submit and view team standups
- 📊 **Attendance**: Mark attendance and view team status
- 🏖️ **Leave Management**: Request and manage leaves
- 🎨 **Beautiful TUI**: Built with Bubble Tea and Lip Gloss
- ⚙️ **Configuration**: Persistent config at `~/.muster/config.yaml`

## Prerequisites

- Go 1.21 or higher
- Muster backend running (see `../backend/README.md`)

## Installation

### Option 1: Build from source

```bash
# Install dependencies
make deps

# Build the binary
make build

# Run the CLI
./bin/muster
```

### Option 2: Install to $GOPATH/bin

```bash
make install
muster
```

### Option 3: Run directly

```bash
make run
```

## Usage

### First Time Setup

When you run the CLI for the first time, you'll be presented with a menu:

1. **Register**: Create a new account and organization
2. **Login**: Sign in to your existing account
3. **Quit**: Exit the application

### Registration

Select "Register" and fill in the following fields:
- **Your Name**: Your full name (2-255 characters)
- **Email**: Your email address
- **Password**: Must be at least 8 characters
- **Organization Name**: Your company/team name (2-255 characters)

Press `Enter` to submit. Your credentials will be saved to `~/.muster/config.yaml`.

### Login

Select "Login" and enter your credentials:
- **Email**: Your registered email
- **Password**: Your password

Press `Enter` to submit.

### Configuration

Configuration is stored at `~/.muster/config.yaml`:

```yaml
user:
  id: 1
  email: user@example.com
  name: John Doe
  role: primary

organization:
  id: 1
  name: Acme Corp
  status: active

auth:
  token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

api:
  base_url: http://localhost:3000
```

**Note**: This file contains sensitive information. Never commit it to version control.

## Navigation

- `↑` / `↓` or `k` / `j`: Navigate menu items
- `Tab` / `Shift+Tab`: Cycle through form fields
- `Enter` / `Space`: Select/Submit
- `Esc` / `Ctrl+C` / `q`: Quit

## Development

### Project Structure

```
cli/
├── cmd/
│   └── muster/
│       └── main.go           # Entry point
├── internal/
│   ├── api/                  # API client
│   │   ├── client.go         # HTTP client wrapper
│   │   └── auth.go           # Auth endpoints
│   ├── auth/                 # Authentication logic
│   │   └── auth.go           # Login/register service
│   ├── config/               # Configuration management
│   │   └── config.go         # Config file I/O
│   └── ui/                   # Bubble Tea UI components
│       ├── styles.go         # Lip Gloss styles
│       ├── menu.go           # Main menu
│       ├── login.go          # Login screen
│       └── register.go       # Registration screen
├── Makefile                  # Build automation
├── go.mod                    # Go module definition
└── README.md                 # This file
```

### Available Commands

```bash
make build         # Build the CLI binary
make run           # Run the CLI directly
make deps          # Install dependencies
make install       # Install binary to $GOPATH/bin
make clean         # Clean build artifacts
make test          # Run tests
make test-coverage # Run tests with coverage
make lint          # Run linters (requires golangci-lint)
make fmt           # Format code
make build-all     # Build for all platforms
make help          # Show help message
```

### Building for Multiple Platforms

```bash
make build-all
```

This creates binaries for:
- macOS (Intel & Apple Silicon)
- Linux (amd64 & arm64)
- Windows (amd64)

Binaries are output to `bin/` directory.

### Adding New Features

1. **API Endpoints**: Add new methods to `internal/api/`
2. **Business Logic**: Add services to `internal/auth/` or create new packages
3. **UI Screens**: Add new Bubble Tea models to `internal/ui/`
4. **Configuration**: Extend `internal/config/config.go` as needed

## Dependencies

- **Bubble Tea**: TUI framework
- **Bubbles**: Pre-built UI components
- **Lip Gloss**: Terminal styling
- **Viper**: Configuration management
- **Resty**: HTTP client

## Troubleshooting

### Connection Errors

If you see connection errors, ensure the backend is running:

```bash
cd ../backend
npm run dev
```

The backend should be running at `http://localhost:3000`.

### Configuration Issues

If you encounter config issues, delete the config file and re-login:

```bash
rm ~/.muster/config.yaml
```

### Build Errors

Ensure you have Go 1.21+ installed:

```bash
go version
```

Install dependencies:

```bash
make deps
```

## Contributing

See the main project README for contribution guidelines.

## License

See the main project LICENSE file.
