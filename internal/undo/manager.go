package undo

import (
	"github.com/coliva/tsk/internal/model"
)

// Command represents an undoable action
type Command interface {
	// Execute performs the command
	Execute(board *model.Board) error
	// Undo reverses the command
	Undo(board *model.Board) error
	// Description returns a human-readable description
	Description() string
}

// Manager handles undo/redo operations
type Manager struct {
	undoStack []Command
	redoStack []Command
	maxSize   int
}

// NewManager creates a new undo manager
func NewManager(maxSize int) *Manager {
	if maxSize <= 0 {
		maxSize = 20
	}
	return &Manager{
		undoStack: make([]Command, 0, maxSize),
		redoStack: make([]Command, 0, maxSize),
		maxSize:   maxSize,
	}
}

// Execute runs a command and adds it to the undo stack
func (m *Manager) Execute(board *model.Board, cmd Command) error {
	if err := cmd.Execute(board); err != nil {
		return err
	}

	// Add to undo stack
	m.undoStack = append(m.undoStack, cmd)
	if len(m.undoStack) > m.maxSize {
		m.undoStack = m.undoStack[1:]
	}

	// Clear redo stack (new action invalidates redo history)
	m.redoStack = m.redoStack[:0]

	return nil
}

// Undo reverses the last command
func (m *Manager) Undo(board *model.Board) (Command, error) {
	if len(m.undoStack) == 0 {
		return nil, nil
	}

	// Pop from undo stack
	cmd := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]

	// Execute undo
	if err := cmd.Undo(board); err != nil {
		// Re-add to undo stack on failure
		m.undoStack = append(m.undoStack, cmd)
		return nil, err
	}

	// Add to redo stack
	m.redoStack = append(m.redoStack, cmd)
	if len(m.redoStack) > m.maxSize {
		m.redoStack = m.redoStack[1:]
	}

	return cmd, nil
}

// Redo re-applies the last undone command
func (m *Manager) Redo(board *model.Board) (Command, error) {
	if len(m.redoStack) == 0 {
		return nil, nil
	}

	// Pop from redo stack
	cmd := m.redoStack[len(m.redoStack)-1]
	m.redoStack = m.redoStack[:len(m.redoStack)-1]

	// Execute the command again
	if err := cmd.Execute(board); err != nil {
		// Re-add to redo stack on failure
		m.redoStack = append(m.redoStack, cmd)
		return nil, err
	}

	// Add back to undo stack
	m.undoStack = append(m.undoStack, cmd)
	if len(m.undoStack) > m.maxSize {
		m.undoStack = m.undoStack[1:]
	}

	return cmd, nil
}

// CanUndo returns true if there are commands to undo
func (m *Manager) CanUndo() bool {
	return len(m.undoStack) > 0
}

// CanRedo returns true if there are commands to redo
func (m *Manager) CanRedo() bool {
	return len(m.redoStack) > 0
}

// UndoCount returns the number of undoable commands
func (m *Manager) UndoCount() int {
	return len(m.undoStack)
}

// RedoCount returns the number of redoable commands
func (m *Manager) RedoCount() int {
	return len(m.redoStack)
}

// Clear removes all commands from both stacks
func (m *Manager) Clear() {
	m.undoStack = m.undoStack[:0]
	m.redoStack = m.redoStack[:0]
}

// LastUndoDescription returns the description of the last undoable command
func (m *Manager) LastUndoDescription() string {
	if len(m.undoStack) == 0 {
		return ""
	}
	return m.undoStack[len(m.undoStack)-1].Description()
}

// LastRedoDescription returns the description of the last redoable command
func (m *Manager) LastRedoDescription() string {
	if len(m.redoStack) == 0 {
		return ""
	}
	return m.redoStack[len(m.redoStack)-1].Description()
}
