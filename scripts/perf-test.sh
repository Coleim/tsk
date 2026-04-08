#!/bin/bash
# Performance test script for tsk
# Generates a board with 500+ tasks to test performance

set -e

TSK_DIR="$HOME/.tsk"
BOARDS_DIR="$TSK_DIR/data/boards"
BOARD_FILE="$BOARDS_DIR/board-perf-test.json"

echo "Creating performance test board with 500+ tasks..."

# Ensure directory exists
mkdir -p "$BOARDS_DIR"

# Generate board JSON
cat > "$BOARD_FILE" << 'HEADER'
{
  "id": "perf-test",
  "name": "Performance Test (500+ tasks)",
  "tasks": [
HEADER

# Generate 600 tasks (200 per status)
task_count=0
for status in "todo" "in_progress" "done"; do
    for i in $(seq 1 200); do
        task_count=$((task_count + 1))
        priority=$((i % 4))  # 0-3 priority rotation
        
        # Add comma for all but last task
        if [ $task_count -lt 600 ]; then
            comma=","
        else
            comma=""
        fi
        
        cat >> "$BOARD_FILE" << TASK
    {
      "id": "task-$task_count",
      "title": "Task $task_count - Performance Test Item",
      "description": "This is a test task for performance testing. Task number $task_count with status $status.",
      "status": "$status",
      "priority": $priority,
      "labels": ["perf-test", "batch-$((i / 50 + 1))"],
      "position": $i,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }$comma
TASK
    done
done

# Close the JSON
cat >> "$BOARD_FILE" << 'FOOTER'
  ],
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
FOOTER

echo "Created board with $task_count tasks at: $BOARD_FILE"
echo ""
echo "To test performance:"
echo "  1. Run: go run ./cmd/tsk or make run"
echo "  2. Press 'b' to open board selector"
echo "  3. Select 'Performance Test (500+ tasks)'"
echo ""
echo "Performance checklist:"
echo "  [ ] App loads within 2 seconds"
echo "  [ ] Navigation (j/k) is responsive (<100ms)"
echo "  [ ] Pane switching (h/l) is smooth"
echo "  [ ] Search responds within 500ms"
echo "  [ ] Scrolling long lists is smooth"
echo "  [ ] Memory usage stays reasonable (<100MB)"
echo ""
echo "To clean up after testing:"
echo "  rm '$BOARD_FILE'"
