## Why

The current UI uses Lipgloss but lacks visual polish and consistency. Task cards blend together, panels lack depth, and the color scheme could have better contrast. Enhancing the visual design will improve usability, make the app feel more professional, and reduce eye strain during extended use.

## What Changes

- Add configurable theme support (dark/light mode)
- Improve task card design with better borders, spacing, and visual hierarchy
- Add subtle gradients and shadows to panels and modals
- Enhance color contrast for accessibility
- Polish status bar, tabs, and modal components
- Add visual indicators for focused/hovered elements
- Improve empty state displays

## Capabilities

### New Capabilities
- `theming`: Theme system supporting dark/light modes with consistent color palettes

### Modified Capabilities
- `tui-core`: Update visual presentation of panels, borders, and spacing for better visual hierarchy

## Impact

- `internal/styles/styles.go`: Major refactor with theme system and enhanced styles
- `internal/ui/*.go`: Update components to use new theme-aware styles
- All visual components will have improved appearance
- No breaking changes to keyboard shortcuts or data structures
