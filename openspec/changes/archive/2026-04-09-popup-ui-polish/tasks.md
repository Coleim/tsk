## 1. Style Fixes

- [x] 1.1 Update PopupSelectedItemStyle() to use softer background (CurrentTheme.Surface instead of Accent)
- [x] 1.2 Update PopupItemStyle() to remove padding or fix background bleeding
- [x] 1.3 Update LabelEditor selection to use arrow-only style (remove TaskSelectedStyle)

## 2. Due Date Popup Conversion

- [x] 2.1 Modify DueDateEditor.View() to remove width/height params and return compact content
- [x] 2.2 Update app.go ModeDueDate rendering to use overlayDialog() wrapper

## 3. Labels Popup Conversion

- [x] 3.1 Modify LabelEditor.View() to remove width/height params and return compact content
- [x] 3.2 Update app.go ModeLabels rendering to use overlayDialog() wrapper

## 4. Verification

- [x] 4.1 Test Due Date popup appearance and keyboard navigation
- [x] 4.2 Test Labels popup appearance and keyboard navigation
- [x] 4.3 Test Edit Task label popup styling (selected and unselected)
- [x] 4.4 Verify search popup still looks correct after style changes
