## Context

The TUI application uses a Catppuccin Mocha-inspired dark theme defined in `internal/styles/styles.go`. The project website at `docs/` uses a GitHub-inspired dark theme in `docs/css/style.css`. These palettes differ significantly:

| Role | TUI (Catppuccin) | Website (GitHub) |
|------|-----------------|------------------|
| Background | `#1e1e2e` | `#0d1117` |
| Surface | `#313244` | `#161b22` |
| Card | `#45475a` | `#21262d` |
| Accent | `#cba6f7` (purple) | `#58a6ff` (blue) |
| Text Primary | `#cdd6f4` | `#f0f6fc` |
| Text Secondary | `#a6adc8` | `#8b949e` |

Building multiple TUI tools with consistent branding requires colors to be importable from a shared location.

## Goals / Non-Goals

**Goals:**
- Create a reusable `pkg/palette` package that can be imported by any Go TUI project
- Website uses the same Catppuccin Mocha palette as the TUI
- Single source of truth for color values
- Brand consistency across TUI and website

**Non-Goals:**
- Changing the TUI color scheme
- Supporting multiple website themes (light mode)
- Publishing palette as a separate Go module (stays in this repo for now)

## Decisions

### 1. Create `pkg/palette/` package
A standalone Go package at `pkg/palette/palette.go` containing:
- Color hex constants (e.g., `Background = "#1e1e2e"`)
- Theme structs with all colors grouped
- Both Dark and Light theme definitions

This package has no dependencies on lipgloss, making it portable.

**Alternative considered**: JSON/YAML config file → rejected because Go constants provide compile-time safety and IDE support.

### 2. Internal styles imports from pkg/palette
Update `internal/styles/styles.go` to import colors from `pkg/palette` instead of hardcoding hex values. This ensures the internal code uses the same source of truth.

### 3. Map website CSS variables to palette values
Update `:root` CSS custom properties to use the same hex values:

```css
:root {
  --bg-primary: #1e1e2e;      /* palette.Background */
  --bg-secondary: #313244;    /* palette.Surface */
  --bg-card: #45475a;         /* palette.Elevated */
  --text-primary: #cdd6f4;    /* palette.TextPrimary */
  --text-secondary: #a6adc8;  /* palette.TextSecondary */
  --text-muted: #6c7086;      /* palette.TextMuted */
  --accent: #74c7ec;          /* palette.Accent (Sapphire) */
  --accent-hover: #8fd4f0;    /* palette.AccentHover */
  --border: #585b70;          /* palette.BorderLight */
  --success: #a6e3a1;         /* palette.Success */
  --warning: #f9e2af;         /* palette.Warning */
  --danger: #f38ba8;          /* palette.Error */
}
```

## Risks / Trade-offs

**[Risk]** Purple accent may reduce contrast on some monitors → Catppuccin is accessibility-tested, mitigated.

**[Risk]** Internal styles must convert string to lipgloss.Color → Minor overhead, acceptable.

## Decision: Accent Color

**Chosen: Catppuccin Sapphire `#74c7ec`** — Bubbly, playful, "bubble charm" vibe.

Hover variant: `#8fd4f0`
