# asana-cli

CLI for managing Asana resources — projects, tasks, portfolios, custom fields, attachments, comments, and dependencies. Includes a Claude Code skill for natural language interaction.

## Quick Start

### macOS / Linux

```bash
curl -sSfL https://raw.githubusercontent.com/danilodrobac/asana-cli/main/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/danilodrobac/asana-cli/main/install.ps1 | iex
```

### Update

```bash
asana-cli update
```

### From source (developers)

Requires Go installed:

```bash
git clone https://github.com/danilodrobac/asana-cli.git
cd asana-cli
make install
```

## Authentication Setup

### 1. Create an OAuth App (one person per team)

1. Go to https://app.asana.com/0/developer-console
2. Create a new app
3. Go to **OAuth** in the sidebar
4. Under **Permission scopes**, toggle **Full permissions**
5. Set the **Redirect URI** to `http://localhost:8931/callback`
6. Copy the **Client ID** and **Client Secret**
7. Share these with your team (they identify the app, not the user)

### 2. Find Your Workspace GID

Open any Asana project in the browser. The URL looks like:
```
https://app.asana.com/0/1234567890/...
                        ^^^^^^^^^^
                        this is your workspace GID
```

### 3. Set Environment Variables

You need three env vars set permanently in your shell. Everyone on the team uses the same Client ID and Secret.

#### macOS — fish shell

Add to `~/.config/fish/config.fish`:

```fish
set -gx ASANA_CLIENT_ID "your-client-id"
set -gx ASANA_CLIENT_SECRET "your-client-secret"
set -gx ASANA_WORKSPACE_ID "your-workspace-gid"
```

Then: `source ~/.config/fish/config.fish`

#### macOS / Linux — bash or zsh

Add to `~/.bashrc` or `~/.zshrc`:

```bash
export ASANA_CLIENT_ID="your-client-id"
export ASANA_CLIENT_SECRET="your-client-secret"
export ASANA_WORKSPACE_ID="your-workspace-gid"
```

Then: `source ~/.bashrc` (or `~/.zshrc`)

#### Windows — PowerShell

Run these commands (sets permanently for your user account):

```powershell
[System.Environment]::SetEnvironmentVariable('ASANA_CLIENT_ID', 'your-client-id', 'User')
[System.Environment]::SetEnvironmentVariable('ASANA_CLIENT_SECRET', 'your-client-secret', 'User')
[System.Environment]::SetEnvironmentVariable('ASANA_WORKSPACE_ID', 'your-workspace-gid', 'User')
```

Then restart your terminal.

### 4. Log In

```bash
asana-cli auth login    # opens browser for OAuth authorization
asana-cli auth status   # verify it worked
```

Each person runs `auth login` separately — you each get your own access tokens.

### Alternative: Personal Access Token

If you don't want OAuth, generate a PAT at https://app.asana.com/0/developer-console and set:

```bash
export ASANA_ACCESS_TOKEN="your-pat"
```

This skips the OAuth flow entirely.

## Usage

```
asana-cli <resource> <verb> [flags]
```

All commands output JSON:
```json
{"success": true, "data": {...}, "message": "Task created successfully"}
```

### Commands

| Resource | Verbs |
|----------|-------|
| `auth` | `login`, `logout`, `status` |
| `project` | `create`, `get`, `list`, `update`, `delete` |
| `task` | `create`, `get`, `list`, `update`, `delete`, `search` |
| `portfolio` | `create`, `get`, `list`, `update`, `delete`, `add-item`, `remove-item` |
| `custom-field` | `create`, `get`, `list`, `update`, `delete` |
| `attachment` | `upload`, `get`, `list`, `delete` |
| `comment` | `create`, `get`, `list`, `update`, `delete` |
| `dependency` | `add`, `remove`, `list` |
| `version` | *(prints version)* |
| `update` | *(self-update to latest release)* |

### Examples

```bash
# Create a project
asana-cli project create --name "Q1 Sprint"

# Create a task in that project
asana-cli task create --name "Design API" --project 1234567890 --assignee user@example.com --due-on 2025-06-15

# Add a comment
asana-cli comment create --task 9876543210 --text "Ready for review"

# Search for tasks
asana-cli task search --query "bug fix"

# Add a dependency
asana-cli dependency add --task 111 --depends-on 222
```

Use `asana-cli <resource> --help` for full flag details.

## Claude Code Integration

After running `install.sh`, the Claude Code skill is active globally. In any Claude Code conversation, you can say things like:

- "Create an Asana task called 'Fix login bug' in project X"
- "List all tasks in the Sprint project"
- "Add a comment to ticket 12345"
- "What Asana tasks are assigned to me?"

Claude will use `asana-cli` commands automatically.

## Development

```bash
make build    # build binary
make test     # run tests
make vet      # run go vet
make install  # build + install binary + install skill
make clean    # remove binary
```
