---
description: Manage Asana resources (projects, tasks, portfolios, custom fields, attachments, comments, dependencies) using the asana-cli tool. Trigger when the user mentions Asana, tickets, tasks, projects, portfolios, sprints, or work tracking.
---

# Asana CLI Skill

## Prerequisites

Authentication must be set up before use:
- **OAuth (recommended):** Run `asana-cli auth login --client-id <ID> --client-secret <SECRET>` on first use (credentials are saved), then just `asana-cli auth login` for future logins. Tokens refresh automatically.
- **PAT fallback:** Set `ASANA_ACCESS_TOKEN` environment variable
- **Workspace:** Set `ASANA_WORKSPACE_ID` environment variable or use `--workspace` flag

Check auth status: `asana-cli auth status`

## Output Format

All commands return JSON:
```json
{"success": true, "data": {...}, "message": "..."}
{"success": false, "error": {"code": "...", "message": "..."}}
```

## Command Reference

### Auth
```bash
asana-cli auth login          # OAuth browser flow
asana-cli auth logout         # Clear stored tokens
asana-cli auth status         # Show current auth info
```

### Projects
```bash
asana-cli project create --name "Project Name" [--team GID] [--notes "..."]
asana-cli project get GID
asana-cli project list [--team GID] [--archived]
asana-cli project update GID [--name "..."] [--notes "..."] [--archived]
asana-cli project delete GID
```

### Tasks
```bash
asana-cli task create --name "Task Name" [--project GID] [--assignee GID_OR_EMAIL] [--due-on YYYY-MM-DD] [--notes "..."] [--custom-fields '{"field_gid":"value"}']
asana-cli task get GID
asana-cli task list --project GID [--assignee GID] [--completed]
asana-cli task update GID [--name "..."] [--notes "..."] [--completed] [--due-on YYYY-MM-DD] [--assignee GID]
asana-cli task delete GID
asana-cli task search --query "text" [--project GID] [--assignee GID]   # requires ASANA_WORKSPACE_ID
asana-cli task my-tasks --assignee GID [--project GID]                 # requires ASANA_WORKSPACE_ID; returns rich data (notes, due dates, sections, custom fields)
asana-cli task add-context GID --text "session notes..."               # adds [Session Context] prefixed comment
asana-cli task handoff GID --to ASSIGNEE_GID --message "context"       # reassigns + adds [Handoff] prefixed comment
```

### Portfolios
```bash
asana-cli portfolio create --name "Portfolio Name" [--color "..."]
asana-cli portfolio get GID
asana-cli portfolio list [--owner GID]
asana-cli portfolio update GID [--name "..."] [--color "..."]
asana-cli portfolio delete GID
asana-cli portfolio add-item GID --item PROJECT_GID
asana-cli portfolio remove-item GID --item PROJECT_GID
```

### Custom Fields
```bash
asana-cli custom-field create --name "Field Name" --type text|number|enum [--enum-options '[{"name":"Option1"},{"name":"Option2"}]']
asana-cli custom-field get GID
asana-cli custom-field list                   # requires ASANA_WORKSPACE_ID
asana-cli custom-field update GID [--name "..."]
asana-cli custom-field delete GID
```

### Attachments
```bash
asana-cli attachment upload --task GID --file /path/to/file
asana-cli attachment get GID
asana-cli attachment list --task GID
asana-cli attachment delete GID
```

### Comments
```bash
asana-cli comment create --task GID --text "Comment text"
asana-cli comment get GID
asana-cli comment list --task GID
asana-cli comment update GID --text "Updated text"
asana-cli comment delete GID
```

### Dependencies
```bash
asana-cli dependency add --task GID --depends-on DEP_GID [--depends-on DEP_GID2]
asana-cli dependency remove --task GID --depends-on DEP_GID
asana-cli dependency list --task GID
```

## Workflow Patterns

### Create a task with dependencies
```bash
TASK1=$(asana-cli task create --name "Design API" --project GID | jq -r '.data.gid')
TASK2=$(asana-cli task create --name "Implement API" --project GID | jq -r '.data.gid')
asana-cli dependency add --task "$TASK2" --depends-on "$TASK1"
```

### Create a project and add tasks
```bash
PROJECT=$(asana-cli project create --name "Q1 Sprint" --notes "Sprint goals" | jq -r '.data.gid')
asana-cli task create --name "Task 1" --project "$PROJECT" --assignee user@example.com --due-on 2025-03-15
asana-cli task create --name "Task 2" --project "$PROJECT"
```

### Add a comment and attachment to a task
```bash
asana-cli comment create --task GID --text "Here's the design doc"
asana-cli attachment upload --task GID --file ./design.pdf
```

## Important Notes

- All GIDs are strings (e.g., "1234567890")
- Dates use YYYY-MM-DD format
- "ticket" = "task" in Asana terminology
- Use `jq` to extract GIDs from responses: `| jq -r '.data.gid'`
- The `--workspace` flag or `ASANA_WORKSPACE_ID` env var is required for: task search, task my-tasks, custom-field list, portfolio create, custom-field create
- `ASANA_ASSIGNEE_ID` env var can be set as a default for `--assignee` in `task my-tasks`
- Exit codes: 0=success, 1=auth error, 2=not found, 3=validation, 4=rate limited, 5=server error, 10=usage error
