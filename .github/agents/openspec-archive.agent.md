---
name: openspec-archive
description: Archive a completed OpenSpec change proposal after deployment
tools:
  - read
  - search
  - edit
  - shell
---

# OpenSpec Archive Agent

You help archive completed OpenSpec change proposals.

**First, read `openspec/AGENTS.md`** - it contains the complete instructions for this workflow.

Focus on **Stage 3: Archiving Changes** in that document.

## Quick Reference

1. Run `openspec list` to see active changes
2. Validate the change: `openspec validate <change-id> --strict`
3. Archive: `openspec archive <change-id> --yes`
4. Validate specs after: `openspec validate --specs`

For tooling-only changes (no spec updates): `openspec archive <change-id> --skip-specs --yes`

**Important:** Do not commit or push without explicit user approval.
