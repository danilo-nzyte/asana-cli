---
description: "Personal work queue powered by Asana. Trigger on: \"what's next\", \"what should I work on\", \"I'm done\", \"hand off to\", \"add a task\", \"what's on my plate\", \"save where I am\", \"show my tasks\", \"what am I working on\", \"create a ticket\", \"assign this to\", \"I'm blocked\", \"park this\", \"switch task\", \"this is urgent\", \"mark this done\", \"continue [task]\", \"pick up where I left off\", \"resume [task]\", \"push this to next week\", \"make this high priority\", \"deprioritise this\""
---

# Work Queue Skill

A personal work queue interface powered by Asana via `asana-cli`. You never open Asana directly — Claude manages it as a byproduct of working with you.

## Prerequisites

- `asana-cli` installed and authenticated (`asana-cli auth status`)
- Config file at `~/.claude/work-queue-vars.yaml` (see First Run below)

## Configuration: `~/.claude/work-queue-vars.yaml`

This file stores your personal Asana context. **Read the entire file, modify in-memory, and write back** (don't append — YAML is whitespace-sensitive).

Expected schema:
```yaml
workspace: "<workspace-gid>"
project: "<default-project-gid>"  # optional — scopes "what's next?" to one project
people:
  me:
    gid: "<your-gid>"
    email: "<your-email>"
  # colleague_name:
  #   gid: "<their-gid>"
  #   name: "Display Name"
custom_fields:
  priority:
    gid: "<gid>"
    values:
      high: "<gid>"
      medium: "<gid>"
      low: "<gid>"
  status:
    gid: "<gid>"
    values:
      not_started: "<gid>"
      on_track: "<gid>"
      at_risk: "<gid>"
      off_track: "<gid>"
      completed: "<gid>"
      on_hold: "<gid>"
sections: {}
  # Populated by Claude as sections are encountered
```

### First Run (vars file doesn't exist)

When `~/.claude/work-queue-vars.yaml` doesn't exist, guide the user through setup:

1. Ask for workspace GID (or detect via `asana-cli project list`)
2. Ask for their Asana user GID/email
3. Ask for default project GID (optional)
4. Run `asana-cli custom-field list` to discover Priority and Status field GIDs
5. Write the vars file with discovered values
6. Proceed with the original request

## Workflows

### "What's next?" / "Show my tasks"

1. Read `~/.claude/work-queue-vars.yaml`
2. Run: `asana-cli task my-tasks --assignee <me.gid> [--project <project>]`
3. Prioritise results using this order:
   - **Priority field**: High > Medium > Low (unset = Medium)
   - **Due date**: overdue first, then soonest
   - **Status**: "At risk" / "Off track" above others
   - **Task age**: older tasks first (tiebreaker)
4. Present top 3–5 tasks: name, section, priority, due date, status, summary from notes
5. Recommend #1 with reasoning that adapts to the actual differentiator (e.g., "due soonest" not "highest priority" when all same priority)
6. On selection, show full notes (structured context)
7. If 100 results returned: "You have many tasks — consider using `--project` to filter."
8. If 0 results: suggest checking project filter, or offer to create a new task

### "Continue X" / "Resume X" / "Pick up where I left off"

1. Search tasks by name: `asana-cli task search --query "X" --assignee <me.gid>`
2. Read full task notes + latest `[Session Context]` comment via `asana-cli comment list --task <GID>`
3. Present "You left off at..." summary from the most recent `[Session Context]` comment
4. Continue working on the task

### During work — creating tasks

When the user asks to add a task or says "[person] needs to...":

1. Ask which section it belongs in (if not obvious)
2. Ask priority (if not stated)
3. Create with structured notes (see template below)
4. Use `asana-cli task create --name "..." --project <GID> --assignee <GID> --notes "..." [--custom-fields '{"<priority_gid>":"<value_gid>"}']`

### "I'm blocked" / "Park this"

1. Update status to "On hold": `asana-cli task update <GID> --custom-fields '{"<status_gid>":"<on_hold_gid>"}'`
   - Note: `task update` doesn't support `--custom-fields` natively. Use `add-context` to log the blocker, and inform the user to update the status field manually in Asana, or use the Asana skill's direct approach.
2. Add context: `asana-cli task add-context <GID> --text "Blocked: <reason>"`
3. Suggest next task

### Priority / date changes

- "Make this urgent" / "high priority": update custom field via task update
- "Push to next week": `asana-cli task update <GID> --due-on <next-monday-date>`
- "Deprioritise": update priority custom field to Low

### "I'm done" / "Mark this done"

Ask: "Is this task complete, or are you wrapping up for now?"

- **Complete**: `asana-cli task update <GID> --completed`
- **Wrapping up**: `asana-cli task add-context <GID> --text "<session summary>"`

Then suggest next task.

### Handoff — "Hand off to [person]"

```bash
asana-cli task handoff <GID> --to <person-gid> --message "<structured handoff>"
```

Handoff message template (structured, not freeform):
```
What was done: ...
What remains: ...
Key files/branches: ...
Blockers/notes: ...
```

### Completion / switching tasks

1. Mark complete OR add session context with where you left off
2. If handing off: use `task handoff`
3. Suggest next task from the queue

## Structured Task Notes Template

```markdown
## Summary
One-line description.

## Context
Why this task exists, what prompted it.

## Git Repo
Repository: [repo-name]
Branch: [branch if relevant]

## Key Files
- `path/to/file.ts` — why relevant

## Acceptance Criteria
- [ ] What "done" looks like

## References
- Links to docs, PRs, related tasks

## Session Log
[Appended via add-context, newest first]
- **YYYY-MM-DD (Person):** What was done. Next steps. Left off at `file.ts:line`.
```

**Not all sections required** — use judgement:
- **Quick tasks** (under 30 min, non-code): Summary only
- **Standard tasks**: Summary + Context + Acceptance Criteria
- **Complex/code tasks**: Full template with Git Repo, Key Files, etc.

## Updating the Sections Cache

When task results include `memberships.section.name` values not yet in the vars file's `sections` map, add them:

```yaml
sections:
  "Section Name": "<section-gid>"
```

Read the vars file, add the new section, write it back.

## Error Recovery

| Exit code | Meaning | Action |
|-----------|---------|--------|
| 1 | Auth error | Tell user: run `asana-cli auth login` and retry |
| 4 | Rate limited | Tell user: wait 60 seconds, then retry |
| 5 | Server error | Tell user: Asana may be down, try again later |

## Limitations

- Asana search API returns max 100 tasks — use `--project` to filter if this is hit
- `task handoff` is not atomic — if comment fails after reassignment, user is warned
- Subtasks are excluded (`is_subtask=false`) by design
- Prioritisation logic is interpreted (non-deterministic) — acceptable for personal use
- No locking — concurrent sessions could both claim the same "next" task
