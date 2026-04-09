# Responsive Layout Specifications

## Layout Modes

### Full Mode (≥80 chars wide)
```
┌──────────────────────────────────────────────────────────────────────────────┐
│ tsk │ Personal Tasks                                   Press ? for help      │
├──────────────────────────────────────────────────────────────────────────────┤
│ ○ TODO (3)    ● IN PROGRESS (2)    ○ DONE (5)                                │
├─────────────────────────────┬────────────────────────────────────────────────┤
│ ▶ Fix login bug          ! │  Fix login bug                                 │
│   Update documentation     │  ────────────────────────────────────────────  │
│   Review PR #123           │                                                 │
│   Write tests              │  Priority: High                                 │
│                            │  Due: Apr 15, 2026                              │
│                            │  Labels: [bug] [urgent]                         │
│                            │                                                 │
│                            │  Description:                                   │
│                            │  The login form throws an error...              │
│       (30 chars min)       │              (remaining)                        │
├─────────────────────────────┴────────────────────────────────────────────────┤
│ 4 tasks │ ↑↓: navigate  Enter: view                                          │
│ Saved                                                                         │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Compact Mode (50-79 chars wide)
```
┌──────────────────────────────────────────────────────────┐
│ tsk │ Personal Tasks                  Press ? for help   │
├──────────────────────────────────────────────────────────┤
│ ○ TODO (3)  ● IN PROGRESS  ○ DONE                        │
├─────────────────────────────┬────────────────────────────┤
│ ▶ Fix login bug          ! │  Fix login bug             │
│   Update documentation     │  ────────────────────────  │
│   Review PR #123           │  Priority: High            │
│   Write tests              │  Due: Apr 15              │
│                            │  [bug] [urgent]            │
│                            │                            │
│       (30 chars fixed)     │     (remaining)            │
├─────────────────────────────┴────────────────────────────┤
│ 4 tasks │ ↑↓ Enter n d                                   │
│ Saved                                                     │
└──────────────────────────────────────────────────────────┘
```

### Single-Panel Mode (<50 chars wide)
```
┌───────────────────────────────────────────┐
│ tsk │ Personal               ? for help   │
├───────────────────────────────────────────┤
│ ○ TODO (3)  ● IN PROGRESS  ○ DONE         │
├───────────────────────────────────────────┤
│ ▶ Fix login bug                        !  │
│   Update documentation                    │
│   Review PR #123                          │
│   Write tests for auth module             │
│   Refactor database layer                 │
│   Add caching                             │
│   Update dependencies                     │
│                                           │
│            (full width)                   │
│                                           │
│          PREVIEW HIDDEN                   │
│     (Enter to view task details)          │
├───────────────────────────────────────────┤
│ 4 tasks │ ↑↓ Enter                        │
│ Saved                                      │
└───────────────────────────────────────────┘
```

### Ultra-Narrow Mode (<35 chars wide)
```
┌─────────────────────────────────┐
│ tsk │ Personal       ? help    │
├─────────────────────────────────┤
│ ● TODO (3)                      │
├─────────────────────────────────┤
│ ▶ Fix login bug              ! │
│   Update docs                   │
│   Review PR #123                │
│   Write tests                   │
│   Refactor DB                   │
│                                 │
│     (PREVIEW HIDDEN)            │
├─────────────────────────────────┤
│ 4 tasks │ ↑↓                    │
│ Saved                           │
└─────────────────────────────────┘
```

## Width Thresholds

| Width | Mode | Task List | Preview |
|-------|------|-----------|---------|
| ≥80 | Full | 30% (min 30) | Remaining |
| 50-79 | Compact | 30 chars fixed | Remaining |
| <50 | Single Panel | Full width | **Hidden** |

## Fixed Elements

- **Header**: Always 1 line, always at top
- **Tabs**: Always 1 line, always below header
- **Status bar**: Always 2 lines, always at bottom
- **Content area**: Fills remaining vertical space
