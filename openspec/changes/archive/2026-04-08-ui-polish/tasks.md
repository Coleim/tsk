## 1. Theme System Foundation

- [x] 1.1 Create Theme struct with all semantic color fields (Background, Surface, Elevated, TextPrimary, TextSecondary, TextMuted, Accent, Success, Warning, Error, PriorityHigh, PriorityMedium, PriorityLow)
- [x] 1.2 Define DarkTheme variable with Catppuccin-inspired colors
- [x] 1.3 Define LightTheme variable with light mode colors
- [x] 1.4 Create CurrentTheme pointer and InitTheme() function
- [x] 1.5 Add TSK_THEME environment variable check in InitTheme()
- [x] 1.6 Update main.go to call InitTheme() on startup

## 2. Refactor Existing Styles

- [x] 2.1 Replace hardcoded ColorPrimary/ColorSecondary etc. with CurrentTheme references
- [x] 2.2 Update TitleStyle to use theme colors
- [x] 2.3 Update BoardNameStyle to use theme colors
- [x] 2.4 Update TabActiveStyle and TabInactiveStyle to use theme colors
- [x] 2.5 Update TaskSelectedStyle and TaskNormalStyle to use theme colors
- [x] 2.6 Update PreviewStyle and PreviewTitleStyle to use theme colors
- [x] 2.7 Update StatusLine1Style and StatusLine2Style to use theme colors
- [x] 2.8 Update ModalStyle and ModalTitleStyle to use theme colors
- [x] 2.9 Update ErrorStyle and SuccessStyle to use theme colors

## 3. Panel Visual Enhancement

- [x] 3.1 Update TaskListStyle with double-line rounded border
- [x] 3.2 Add panel title rendering for task list ("Tasks")
- [x] 3.3 Update PreviewStyle with doubled-line rounded border
- [x] 3.4 Add panel title rendering for preview ("Preview")
- [x] 3.5 Apply surface background color to both panels
- [x] 3.6 Ensure consistent padding (1 vertical, 2 horizontal)

## 4. Tab Visual Polish

- [x] 4.1 Create PillTabStyle for active tabs with filled background
- [x] 4.2 Update active tab to use status-specific colors (ToDo=primary, InProgress=warning, Done=success)
- [x] 4.3 Update inactive tab styling with muted text, no background
- [x] 4.4 Add proper spacing between tabs

## 5. Task Card Enhancement

- [x] 5.1 Create SelectedTaskCardStyle with elevated background and priority accent border
- [x] 5.2 Create UnselectedTaskCardStyle with surface background
- [x] 5.3 Add left border accent colored by task priority
- [x] 5.4 Add 1-unit vertical margin between cards
- [x] 5.5 Update task card rendering to use new styles

## 6. Modal Polish

- [x] 6.1 Update ModalStyle with elevated background color
- [x] 6.2 Increase border thickness for modals
- [x] 6.3 Ensure modal title is bold with accent color
- [x] 6.4 Apply elevated styling to all dialog types (edit, create, filter, help)

## 7. Status Bar Polish

- [x] 7.1 Add top border separator line to status bar
- [x] 7.2 Style line 1 with appropriate color based on message type
- [x] 7.3 Style line 2 with secondary text color
- [x] 7.4 Add subtle background color to status bar area

## 8. Empty States

- [x] 8.1 Create EmptyStateStyle with centered, muted text
- [x] 8.2 Update empty pane message: "No tasks • Press 'n' to create"
- [x] 8.3 Update empty search results message styling
- [x] 8.4 Apply empty state styling to all empty views

## 9. Testing and Documentation

- [x] 9.1 Test dark theme in multiple terminal emulators (iTerm2, Terminal.app, VS Code)
- [x] 9.2 Test light theme in multiple terminal emulators
- [x] 9.3 Verify color contrast meets accessibility standards
- [x] 9.4 Update README with theme documentation (TSK_THEME env var)
- [ ] 9.5 Add screenshots of both themes to documentation
