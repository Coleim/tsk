## Context

The tsk terminal task manager uses Lipgloss for styling but lacks visual polish. Current styles use basic colors without theme structure, making it hard to maintain consistency or add features like dark/light mode switching. The UI works but feels utilitarian rather than polished.

**Current State:**
- Basic color palette defined as constants in `internal/styles/styles.go`
- Styles scattered without clear organization
- No theme system - colors hardcoded
- Panels have minimal visual hierarchy
- Task cards lack visual distinction

## Goals / Non-Goals

**Goals:**
- Create a theme system supporting dark/light modes
- Improve visual hierarchy with better spacing, borders, and depth
- Enhance color contrast for accessibility (WCAG AA)
- Polish all major components (task cards, panels, modals, tabs, status bar)
- Make styles maintainable and consistent

**Non-Goals:**
- Custom user themes (only built-in dark/light)
- Animations beyond existing spinner
- Terminal image support
- Configuration file for themes (env var only)

## Decisions

### Decision: Theme Structure

Use a `Theme` struct containing all semantic colors, then create `DarkTheme` and `LightTheme` instances. Components reference `CurrentTheme` rather than direct colors.

```go
type Theme struct {
    // Backgrounds
    Background   lipgloss.Color
    Surface      lipgloss.Color
    Elevated     lipgloss.Color
    
    // Text
    TextPrimary   lipgloss.Color
    TextSecondary lipgloss.Color
    TextMuted     lipgloss.Color
    
    // Semantic
    Accent  lipgloss.Color
    Success lipgloss.Color
    Warning lipgloss.Color
    Error   lipgloss.Color
    
    // Priority
    PriorityHigh   lipgloss.Color
    PriorityMedium lipgloss.Color
    PriorityLow    lipgloss.Color
    
    // etc.
}
```

**Rationale**: A struct allows type safety and IDE completion. Global `CurrentTheme` pointer makes switching instant.

### Decision: Theme Selection

Use `TSK_THEME` environment variable with values `dark` (default) or `light`. No runtime switching needed.

**Rationale**: Simple, no persistence needed, follows terminal app conventions.

### Decision: Visual Hierarchy

Add three elevation levels using background colors:
1. **Background**: Base terminal color
2. **Surface**: Panels, cards (subtle lift)
3. **Elevated**: Modals, popups (more prominent)

**Rationale**: Creates depth without relying on shadows which don't exist in terminals.

### Decision: Component Polish

| Component | Enhancement |
|-----------|-------------|
| Task cards | Border on selected, padding, priority accent stripe |
| Panels | Double-line borders, titles, proper spacing |
| Tabs | Pill-style with background fill for active |
| Modals | Elevated background, thicker borders |
| Status bar | Clear separation line, icon indicators |
| Empty states | Centered, muted text with helpful hints |

### Decision: Color Palette

Dark theme (default):
- Background: #1e1e2e (Catppuccin-inspired)
- Surface: #313244
- Elevated: #45475a
- Accent: #cba6f7 (Purple/pink)
- Text: #cdd6f4 (Soft white)

Light theme:
- Background: #eff1f5
- Surface: #e6e9ef
- Elevated: #dce0e8
- Accent: #8839ef
- Text: #4c4f69

**Rationale**: Modern, accessible colors inspired by popular terminal themes.

## Risks / Trade-offs

**[Risk]** Terminal color support varies → Use 256-color codes (Color("213")) which work on most modern terminals

**[Risk]** Light theme less tested → Document as "experimental" initially

**[Trade-off]** Fixed themes vs customization → Simplicity wins; users wanting custom themes can edit source

**[Trade-off]** Elevated backgrounds may look odd on some terminals → Provide option to disable via `TSK_FLAT=1`
