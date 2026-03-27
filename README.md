# Muster CLI

> A beautiful terminal user interface (TUI) for Muster - a team productivity tool for daily standups, attendance tracking, leave management, and project organization.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)

---

## Features

- 📝 **Standups**: Submit and view team standups (scoped to projects)
- 📊 **Attendance**: Check in/out with time tracking and view team status
- 🏖️ **Leave Management**: Request (sick, casual, paid, unpaid, WFH), review, and cancel leaves
- 📂 **Projects**: Create projects, manage members, configure per-project settings
- 👥 **Team Management**: Invite members, manage roles
- 🔑 **Password Management**: Change your own password or reset a member's password (admin)
- 🔐 **Authentication**: Login, registration, and org invitations
- 🎨 **Beautiful TUI**: Built with Bubble Tea and Lip Gloss
- 🔒 **Role-based commands**: CLI only shows commands you have access to

## Installation

### Homebrew (macOS & Linux)

```bash
brew install beamlabco/tap/muster
```

### Build from source

```bash
git clone https://github.com/beamlabco/muster.git
cd muster/cli
make build
./bin/muster
```

## Usage

Launch the CLI:

```bash
muster
```

### Commands

Type `/` to see available commands. Commands are filtered by your role — admins see all commands, members see only what they have access to.

#### All Users

| Command | Description |
|---------|-------------|
| `/login` | Sign in to your account |
| `/signup` | Create a new account |
| `/join` | Join an organization via invitation |
| `/standup` | Submit your daily standup |
| `/standup today` | View team standups for today |
| `/standup history` | View your standup history |
| `/checkin` | Check in for the day |
| `/checkout` | Check out for the day |
| `/attendance` | Mark your attendance |
| `/attendance today` | View team attendance for today |
| `/attendance history` | View your attendance history |
| `/leave` | Request a leave (sick, casual, paid, unpaid, wfh) |
| `/leave list` | View team leaves |
| `/leave cancel` | Cancel a pending leave |
| `/project` | View your projects |
| `/team` | View team members and roles |
| `/password` | Change your password |
| `/whoami` | Show current user info |
| `/help` | Show available commands |
| `/quit` | Exit the application |

#### Primary (Admin) Only

| Command | Description |
|---------|-------------|
| `/project create` | Create a new project |
| `/project settings` | Configure project settings (name, summary time, timezone, Discord webhook) |
| `/project members` | Add or remove project members |
| `/role` | Update a user's role |
| `/invite` | Invite a team member |
| `/settings` | Configure organization settings |
| `/reset-password` | Reset a member's password |

#### Primary + Manager

| Command | Description |
|---------|-------------|
| `/leave review` | Review and approve/reject pending leaves |

### Navigation

- `Tab`: Accept autocomplete suggestion
- `↑` / `↓`: Navigate suggestions
- `←` / `→`: Cycle options (e.g., leave type)
- `Enter`: Execute command / submit form
- `Esc`: Go back / cancel
- `Ctrl+C`: Quit

### Configuration

Config is stored at `~/.muster/config.yaml`:

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
  base_url: https://api.muster.stagify.xyz
```

## Development

### Make commands

```bash
make build         # Build the CLI binary
make run           # Run the CLI directly
make deps          # Install dependencies
make install       # Install binary to /usr/local/bin
make clean         # Clean build artifacts
make test          # Run tests
make test-coverage # Run tests with coverage
make lint          # Run linters (requires golangci-lint)
make fmt           # Format code
make build-all     # Build for all platforms
```

### Releasing

Releases are automated via GitHub Actions and [GoReleaser](https://goreleaser.com/). Pushing a version tag triggers a build that publishes binaries to GitHub Releases and updates the Homebrew tap.

```bash
# Create and push a release tag
git tag v0.2.0
git push origin v0.2.0
```

To test a release locally without publishing:

```bash
goreleaser release --snapshot --clean
```

Binaries are built for:
- macOS (Intel & Apple Silicon)
- Linux (amd64 & arm64)
- Windows (amd64)

## License

See the main project LICENSE file.
