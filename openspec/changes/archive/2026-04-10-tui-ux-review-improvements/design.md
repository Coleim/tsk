## Context

The `tsk` TUI application has a well-polished main board view with good visual hierarchy. However, secondary screens (full-screen dialogs and popup overlays) feel sparse and visually inconsistent. The current implementation uses Lipgloss with a Catppuccin-inspired theme, but the styling is not uniformly applied across all views.

**Current State:**
- Main board view: Two-panel layout with task list and preview, tabs, status bar - all well-styled
- Filter view: Full-screen but content is sparse, uses basic checkbox list without visual grouping
- Edit task view: Full-screen form with three fields, minimal visual distinction between sections
- Board selector: Full-screen list without visual cards or metadata display
- Popup overlays (search, labels, due date): Generally good but inconsistent widths and padding

**Constraints:**
- Must work in terminal environments with limited Unicode/styling support
- Should degrade gracefully in terminals without true color support
- Must maintain keyboard-first navigation paradigm
- Cannot introduce new dependencies beyond existing Lipgloss/Bubbletea

## Goals / Non-Goals

**Goals:**
- Create visual consistency between main view and secondary screens
- Reduce sparse/empty feeling in full-screen dialogs through visual grouping
- Establish reusable visual patterns (section cards, dividers, form styles)
- Improve visual hierarchy to guide user attention
- Maintain accessibility standards for text contrast

**Non-Goals:**
- Adding mouse interaction support
- Changing keyboard shortcuts or navigation patterns
- Adding animations (terminal animation support is inconsistent)
- Redesigning the main board view (it's already well-polished)
- Supporting additional color themes beyond dark/light

## Decisions

### Decision 1: Use Section Cards for Visual Grouping
Instead of flat lists in full-screen dialogs, wrap related content in subtle card containers with borders.

**Rationale**: Cards create visual hierarchy without requiring background colors that might not render well in all terminals. The existing `lipgloss.RoundedBorder()` provides a modern look.

**Alternatives Considered:**
- Background shading for sections: Rejected - colors may not render consistently across terminals
- Horizontal dividers only: Rejected - less visually distinct than bordered cards

### Decision 2: Standardize Dialog Layout with Header/Content/Footer Pattern
All full-screen dialogs will follow a consistent three-zone layout:
1. **Header**: Modal title with accent styling
2. **Content**: Main interaction area with section cards
3. **Footer**: Keyboard hints with consistent styling

**Rationale**: Users learn one layout pattern that applies everywhere, reducing cognitive load.

### Decision 3: Use Arrow Indicators + Subtle Background for Focus
Selected/focused items will use the `▶` arrow indicator plus a subtle background color shift, not heavy borders.

**Rationale**: This is already partially implemented (TaskSelectedStyle) and works well. Applying consistently reduces visual noise while maintaining clear focus indication.

### Decision 4: Adaptive Popup Widths
- Simple input popups (new task): 50 characters
- List popups (search, labels): 50-60 characters based on content
- Date picker: 60 characters (accommodates formatted dates)

**Rationale**: Fixed widths prevent jarring resizes while ensuring content fits. These sizes work well at minimum terminal widths.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Visual changes may look different across terminal emulators | Test on iTerm2, Terminal.app, and common Linux terminals; use safe border/color choices |
| More visual elements could slow down rendering | Use Lipgloss efficiently; measure render performance; these changes are styling-only |
| Changes could break existing test assertions | Update visual snapshot tests as part of implementation |
| Users may prefer the current sparse style | Changes are additive visual polish, not restructuring; core interactions unchanged |

## Migration Plan

No migration needed - these are purely visual changes to existing views. Implementation can be done incrementally:

1. Add new styles to `styles.go` 
2. Update Filter view (highest visual improvement potential)
3. Update Edit task view
4. Update Board selector view
5. Polish popup overlays for consistency
