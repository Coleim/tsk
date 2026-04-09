## 1. Layout Constants

- [x] 1.1 Define width threshold constants (MinTaskListWidth=30, SinglePanelThreshold=50)

## 2. Width Calculation

- [x] 2.1 Update renderMainView() to calculate taskListWidth with minimum of 30 chars
- [x] 2.2 Make previewWidth take remaining space (width - taskListWidth)
- [x] 2.3 Implement single-panel mode: hide preview entirely when width < 50
- [x] 2.4 Conditionally render preview panel based on width threshold

## 3. Height Stability

- [x] 3.1 Verify header always renders at top regardless of window height
- [x] 3.2 Verify status bar always renders at bottom regardless of window height
- [x] 3.3 Handle minimum height gracefully (content shrinks, bars stay fixed)

## 4. Testing

- [x] 4.1 Test full mode layout (≥80 chars)
- [x] 4.2 Test compact mode layout (50-79 chars)
- [x] 4.3 Test single-panel mode (<50 chars) - verify preview hidden
- [x] 4.4 Verify no visual artifacts during resize transitions
