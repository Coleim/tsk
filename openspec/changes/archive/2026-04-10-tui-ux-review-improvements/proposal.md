## Why

While the home screen (main board view) has solid UX with good visual hierarchy, the secondary screens — full-screen dialogs and popup overlays — suffer from sparse layouts, inconsistent visual density, and weak visual hierarchy. This creates a jarring user experience when navigating between the polished home screen and the less refined secondary views.

## What Changes

**Full-Screen Dialog Improvements:**
- **Filter View**: Add visual grouping with section borders/cards, improve spacing to reduce empty feel, add filter preview summary with better prominence
- **Edit Task View**: Add visual sections with subtle backgrounds, improve field focus indicators, add character count hints, enhance label autocomplete popup styling
- **Board Selector View**: Add board cards with rounded borders and subtle backgrounds, show board metadata (task counts, last accessed) with better visual hierarchy, improve creation/rename form styling

**Visual Hierarchy Enhancements:**
- Consistent section headers with accent underlines across all dialogs
- Better visual separation between form groups using subtle dividers or spacing
- Prominent keyboard hint bars with consistent styling at dialog bottoms
- Active/focused field indicators with subtle glow or border highlight

**Popup Overlay Refinements:**
- Standardize popup widths based on content type (narrow for simple inputs, wider for lists)
- Add subtle shadow/backdrop effect to improve popup prominence
- Ensure consistent padding and border radius across all popups

**Accessibility & Polish:**
- Review text contrast ratios in both light and dark themes
- Add subtle animations for dialog/popup transitions (optional, can degrade gracefully)
- Improve empty state illustrations/messaging in sparse views

## Capabilities

### New Capabilities
- `visual-components`: Reusable styled components (section cards, dividers, enhanced form fields) for consistent visual language across all screens

### Modified Capabilities
- `tui-core`: Updated visual presentation of full-screen dialogs and popup overlays for improved density and hierarchy

## Impact

- `internal/styles/styles.go`: New section/card styles, enhanced form field styles, improved visual hierarchy utilities
- `internal/ui/filter.go`: Refactor View() with sectioned layout, visual grouping
- `internal/ui/edit.go`: Refactor View() with form sections, improved field styling
- `internal/ui/board_selector.go`: Refactor View() with board cards, metadata display
- `internal/ui/detail.go`: Minor polish to section spacing and hierarchy
- `internal/ui/labels.go`: Minor popup styling improvements
- `internal/ui/duedate.go`: Minor popup styling improvements
- `internal/ui/search.go`: Minor popup styling improvements
- `internal/ui/modal.go`: Consistent popup styling foundation
