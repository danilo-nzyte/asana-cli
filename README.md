# asana-cli

CLI for managing Asana resources — projects, tasks, portfolios, custom fields, attachments, comments, and dependencies. Includes a Claude Code skill for natural language interaction.

## Quick Start

### macOS / Linux

```bash
curl -sSfL https://raw.githubusercontent.com/danilo-nzyte/asana-cli/main/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/danilo-nzyte/asana-cli/main/install.ps1 | iex
```

### Update

```bash
asana-cli update
```

### From source (developers)

Requires Go installed:

```bash
git clone https://github.com/danilo-nzyte/asana-cli.git
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

### 2. Log In

Run `auth login` with your OAuth credentials. This is a one-time setup — credentials are saved to `~/.config/asana-cli/config.json` and tokens refresh automatically.

```bash
asana-cli auth login --client-id "your-client-id" --client-secret "your-client-secret"
```

This opens your browser for Asana authorization. Each team member runs this with the same client ID/secret but authorizes with their own Asana account.

Verify it worked:

```bash
asana-cli auth status
```

### 3. Set Workspace ID

The workspace ID is needed for some commands (task search, custom field list, portfolio create). Set it in your shell config:

#### macOS — fish shell

Add to `~/.config/fish/config.fish`:

```fish
set -gx ASANA_WORKSPACE_ID "your-workspace-gid"
```

#### macOS / Linux — bash or zsh

Add to `~/.bashrc` or `~/.zshrc`:

```bash
export ASANA_WORKSPACE_ID="your-workspace-gid"
```

#### Windows — PowerShell

```powershell
[System.Environment]::SetEnvironmentVariable('ASANA_WORKSPACE_ID', 'your-workspace-gid', 'User')
```

Then restart your terminal.

To find your workspace GID, open any Asana project in the browser. The URL looks like:
```
https://app.asana.com/0/1234567890/...
                        ^^^^^^^^^^
                        this is your workspace GID
```

### Alternative: Personal Access Token

If you don't want OAuth, generate a PAT at https://app.asana.com/0/developer-console and set:

```bash
export ASANA_ACCESS_TOKEN="your-pat"
```

This skips the OAuth flow entirely.

### Alternative: Environment Variables for OAuth

You can also set OAuth credentials as environment variables instead of using `--client-id`/`--client-secret` flags. The CLI checks the config file first, then falls back to env vars.

```bash
export ASANA_CLIENT_ID="your-client-id"
export ASANA_CLIENT_SECRET="your-client-secret"
```

Note: credentials provided via env vars are automatically saved to the config file on first `auth login`, so token refresh works even if the env vars are not set in future sessions.

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
| `section` | `create`, `get`, `list`, `update`, `delete`, `add-task` |
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
