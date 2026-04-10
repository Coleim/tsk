## 1. Foundation Styles

- [x] 1.1 Add SectionCardStyle() to styles.go with rounded border and BorderLight color
- [x] 1.2 Add SectionCardTitleStyle() with accent underline effect
- [x] 1.3 Add FormFieldLabelStyle() for standard form labels (muted color)
- [x] 1.4 Add FormFieldActiveLabelStyle() with arrow indicator and primary color
- [x] 1.5 Add KeyboardHintBarStyle() for consistent footer hints
- [x] 1.6 Add CheckboxCheckedStyle() with success-colored checkmark
- [x] 1.7 Add DialogSeparatorStyle() for thin horizontal dividers

## 2. Filter Dialog Polish

- [x] 2.1 Wrap priority section in SectionCard with "Priority" title
- [x] 2.2 Wrap labels section in SectionCard with "Labels" title
- [x] 2.3 Update checkbox rendering to use CheckboxCheckedStyle
- [x] 2.4 Add active filters summary at top when filters are selected
- [x] 2.5 Add separator line above keyboard hint bar
- [x] 2.6 Adjust overall spacing for better visual balance

## 3. Edit Task Dialog Polish

- [x] 3.1 Apply FormFieldActiveLabelStyle to focused field labels
- [x] 3.2 Apply FormFieldLabelStyle to unfocused field labels
- [x] 3.3 Add consistent spacing between form rows
- [x] 3.4 Add separator line above keyboard hint bar
- [x] 3.5 Improve label autocomplete popup styling with consistent background

## 4. Board Selector Dialog Polish

- [x] 4.1 Wrap each board entry in SectionCard-style container
- [x] 4.2 Add task count display with secondary text styling
- [x] 4.3 Highlight current board card with subtle surface background
- [x] 4.4 Apply consistent arrow indicator for selected board
- [x] 4.5 Add separator line above keyboard hint bar
- [x] 4.6 Improve create/rename form input styling

## 5. Popup Overlay Consistency

- [x] 5.1 Standardize search popup width to 50 characters
- [x] 5.2 Standardize labels popup width to 50 characters
- [x] 5.3 Verify due date popup width is 60 characters
- [x] 5.4 Ensure all popups use consistent border styling (double border, accent color)
- [x] 5.5 Verify consistent internal padding across all popups

## 6. Detail View Minor Polish

- [x] 6.1 Add subtle section dividers between metadata and description
- [x] 6.2 Improve spacing between detail fields
- [x] 6.3 Apply separator line above keyboard hint bar

## 7. Testing and Verification

- [x] 7.1 Run application and verify Filter dialog visual improvements
- [x] 7.2 Run application and verify Edit Task dialog visual improvements
- [x] 7.3 Run application and verify Board Selector visual improvements
- [x] 7.4 Run application and verify popup overlays are consistent
- [x] 7.5 Test in light theme mode (TSK_THEME=light)
- [x] 7.6 Verify all existing tests pass
