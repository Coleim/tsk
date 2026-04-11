## Why

The TUI and website currently use different color palettes (TUI: Catppuccin-inspired, website: GitHub-inspired). This visual inconsistency weakens brand identity. Additionally, building multiple TUI tools with consistent branding requires a shared color palette that can be imported across projects.

## What Changes

- Create a standalone color palette package at `pkg/palette/` that can be imported by any Go TUI project
- Define a canonical color palette in a shared reference document
- Update website CSS variables to match TUI colors (Catppuccin Mocha dark theme)
- Ensure accent, semantic, and background colors are consistent:
  - Background: `#1e1e2e` (TUI) → website
  - Surface: `#313244` → website cards
  - Accent: `#cba6f7` (purple) → website links/CTAs
  - Success: `#a6e3a1`, Warning: `#f9e2af`, Error: `#f38ba8`
  - Text primary: `#cdd6f4`, secondary: `#a6adc8`, muted: `#6c7086`

**Decision**: The palette lives in `pkg/palette/` within this repo. Other projects can import `github.com/coliva/tsk/pkg/palette`. If maintaining multiple TUI tools reveals friction, extract to a separate `tui-palette` repo later—the API is small enough for trivial migration.

## Capabilities

### New Capabilities
- `color-palette`: Reusable color palette package for Go TUI tools with canonical color definitions

### Modified Capabilities
- `static-website`: Website CSS uses unified color palette

## Impact

- `pkg/palette/palette.go`: New standalone package with color definitions (importable as `github.com/coliva/tsk/pkg/palette`)
- `docs/css/style.css`: Update CSS custom properties to use Catppuccin Mocha colors
- `internal/styles/styles.go`: Import colors from `pkg/palette` instead of hardcoding
- Documentation: Add palette reference for future development
