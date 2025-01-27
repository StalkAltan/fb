package main

import (
	"sync"
)

type SelectionManager struct {
	selectedPaths map[string]bool
	mu            sync.RWMutex
}

func NewSelectionManager() *SelectionManager {
	return &SelectionManager{
		selectedPaths: make(map[string]bool),
	}
}

func (sm *SelectionManager) Add(path string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.selectedPaths[path] = true
}

func (sm *SelectionManager) Remove(path string) {
	if path == "" {
		return // Early return for empty path
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if map exists
	if sm.selectedPaths == nil {
		return
	}

	// Check existence before delete
	if _, exists := sm.selectedPaths[path]; exists {
		delete(sm.selectedPaths, path)
	}
}

func (sm *SelectionManager) Toggle(path string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.selectedPaths[path] {
		delete(sm.selectedPaths, path)
	} else {
		sm.selectedPaths[path] = true
	}
}

func (sm *SelectionManager) IsSelected(path string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.selectedPaths[path]
}

func (sm *SelectionManager) GetSelectedPaths() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	paths := make([]string, 0, len(sm.selectedPaths))
	for path := range sm.selectedPaths {
		paths = append(paths, path)
	}
	return paths
}
